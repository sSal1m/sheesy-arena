"""Generate source-backed system inventories; never reads event.db or credentials.

Run from any directory: python tools/generate_system_map.py
Uses only Python's standard library. The narrative map is maintained separately.
"""

import json
from pathlib import Path
import re
import subprocess


ROOT = Path(__file__).resolve().parents[1]
OUT = ROOT / "docs" / "harita"


def read(path):
    return (ROOT / path).read_text(encoding="utf-8-sig")


def balanced_end(text, start, opening="{", closing="}"):
    """Find a Go block's end while ignoring comments, quoted and raw strings."""
    depth, index = 0, start
    while index < len(text):
        if text.startswith("//", index):
            end = text.find("\n", index)
            index = len(text) if end < 0 else end
            continue
        if text.startswith("/*", index):
            end = text.find("*/", index + 2)
            index = len(text) if end < 0 else end + 2
            continue
        char = text[index]
        if char in ('"', "'", "`"):
            quote = char
            index += 1
            while index < len(text):
                if quote != "`" and text[index] == "\\":
                    index += 2
                elif text[index] == quote:
                    index += 1
                    break
                else:
                    index += 1
            continue
        if char == opening:
            depth += 1
        elif char == closing:
            depth -= 1
            if depth == 0:
                return index + 1
        index += 1
    raise ValueError("Unbalanced source block")


def functions(path, source):
    result = []
    pattern = r"^func\s+(?:\(([^\n)]*)\)\s+)?(\w+)(?:\[[^\n]*?\])?\("
    for match in re.finditer(pattern, source, re.M):
        start = source.find("{", match.end())
        end = balanced_end(source, start)
        result.append({
            "name": match[2], "receiver": match[1] or "", "file": path,
            "line": source.count("\n", 0, match.start()) + 1,
            "body": source[start:end],
        })
    return result


def link(path, line=None, label=None):
    anchor = "#L" + str(line) if line else ""
    return f"[{label or path}](../../{path}{anchor})"


def cell(value):
    return str(value).replace("|", "\\|").replace("\n", " ")


def code(value):
    if "`" in value:
        return "`` " + value + " ``"
    return "`" + value + "`"


def table(headers, rows):
    return "\n".join([
        "| " + " | ".join(headers) + " |",
        "| " + " | ".join("---" for _ in headers) + " |",
        *("| " + " | ".join(cell(v) for v in row) + " |" for row in rows),
    ]) + "\n"


def write(name, content):
    (OUT / name).write_text(content.rstrip() + "\n", encoding="utf-8")


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    tracked = subprocess.check_output(
        ["git", "ls-files"], cwd=ROOT, text=True, encoding="utf-8"
    ).splitlines()
    gofiles = sorted(p for p in tracked if p.endswith(".go"))
    production = [p for p in gofiles if not p.endswith("_test.go")]
    sources = {p: read(p) for p in production}
    all_functions = [f for p, source in sources.items() for f in functions(p, source)]
    web_functions = {f["name"]: f for f in all_functions if f["file"].startswith("web/")}
    routes = []
    for match in re.finditer(r'mux.HandleFunc\("(\w+) ([^"\n]+)", web\.(\w+)\)', sources["web/web.go"]):
        method, path, handler = match.groups()
        fn = web_functions[handler]
        routes.append({
            "method": method, "path": path, "handler": handler,
            "file": fn["file"], "line": fn["line"],
            "direct_admin_check": "web.userIsAdmin(w, r)" in fn["body"],
            "conditional_fta_auth": "web.userIsAdmin(w, r)" in fn["body"] and handler in (
                "fieldMonitorDisplayHandler", "fmsFieldMonitorDisplayHandler", "fieldMonitorDisplayWebsocketHandler"),
            "templates": sorted(set(re.findall(r'"(templates/[^"\n]+)"', fn["body"]))),
        })
    route_rows = []
    for route in routes:
        auth = "Admin kontrolü" if route["direct_admin_check"] else "Doğrudan admin kontrolü yok"
        if route["conditional_fta_auth"]:
            auth = "Yalnızca fta=true ise admin kontrolü"
        route_rows.append([
            route["method"], "`" + route["path"] + "`",
            link(route["file"], route["line"], route["handler"]), auth,
        ])
    write("HTTP-ROTALARI.md", "# HTTP rota envanteri\n\n"
          "Kaynak: `web/web.go` içindeki açık HandleFunc kayıtları. "
          "Admin sütunu ilgili handler gövdesindeki doğrudan `userIsAdmin` çağrısını gösterir; "
          "bağımsız bir güvenlik denetimi değildir. Admin parolası boşsa bu kontrol de erişimi açık bırakır.\n\n"
          f"Toplam **{len(routes)}** açık rota; **{sum('/websocket' in r['path'] for r in routes)}** WebSocket rotası. "
          "Ek olarak `/static/` dosya sunucusu ve `/` yönlendiricisi ServeWebInterface içinde kaydedilir. "
          "GET kalıpları Go ServeMux davranışı gereği HEAD isteklerini de karşılayabilir.\n\n"
          + table(["Yöntem", "Yol", "Handler / kaynak", "Erişim kontrolü"], route_rows))

    websocket_entries = []
    for route in routes:
        if not route["path"].endswith("/websocket"):
            continue
        fn = web_functions[route["handler"]]
        body = fn["body"]
        subscriptions = []
        for match in re.finditer(r"\bHandleNotifiers\(", body):
            start = match.end() - 1
            end = balanced_end(body, start, "(", ")")
            subscriptions.extend(re.findall(r"web\.arena\.([\w.]+(?:\(\))?)|\b(display\.Notifier)", body[start:end]))
        subs = [a or b for a, b in subscriptions]
        commands = []
        # Command switch cases are two tabs deep; nested value switches are deeper.
        for line in re.findall(r"^\t\tcase ([^\n]+):", body, re.M):
            commands.extend(re.findall(r'"([^"\n]+)"', line))
        commands.extend(re.findall(r'\b(?:command|messageType) == "([^"\n]+)"', body))
        direct_messages = sorted(set(re.findall(r'(?:ws\.Write|writeWebsocketMessage)\((?:ws,\s*)?"([^"\n]+)"', body)))
        websocket_entries.append({**route, "subscriptions": subs,
                                  "commands": list(dict.fromkeys(commands)),
                                  "direct_messages": direct_messages})
    notifier_entries = []
    for match in re.finditer(r'arena\.(\w+) = websocket.NewNotifier\(\s*"([^"\n]+)",\s*([^\s,)]+)', sources["field/arena_notifiers.go"]):
        notifier_entries.append({"field": match[1], "type": match[2], "producer": match[3]})
    ws_content = "# WebSocket sözleşmeleri\n\n"
    ws_content += "Üretilen envanter; veri alanlarının anlamı ana sistem haritasında açıklanır. "
    ws_content += "Çıkarım, handler içindeki HandleNotifiers çağrıları ve üst seviye komut dalları üzerinden yapılır. "
    ws_content += "`display.Notifier` dinamik ekran URL yapılandırmasıdır. `reload`, `playSound` gibi üreticisi nil olaylar bağlantı başında gönderilmez.\n\n"
    ws_content += "## Mesaj zarfı\n\n```json\n{\"type\": \"startMatch\", \"data\": {\"MuteMatchSounds\": false}}\n```\n\n"
    ws_content += "Alan isimleri veri gövdelerindeki Go struct adlarına göre çoğunlukla PascalCase biçimindedir. "
    ws_content += "`ping` ve `error` altyapı mesajlarıdır. Her bağlantı kendi notifier kanallarına abone olur; kanallar arası toplam sıra garantisi yoktur.\n\n"
    ws_content += "## Arena olayları\n\n" + table(
        ["Arena alanı", "JSON type", "Veri üreticisi"],
        [["`" + n["field"] + "`", "`" + n["type"] + "`", "`" + n["producer"] + "`"] for n in notifier_entries],
    )
    for entry in websocket_entries:
        ws_content += "\n## `" + entry["path"] + "`\n\n"
        ws_content += link(entry["file"], entry["line"], entry["handler"]) + "\n\n"
        auth = "yalnızca fta=true ise" if entry["conditional_fta_auth"] else ("var" if entry["direct_admin_check"] else "doğrudan yok")
        ws_content += "- Admin kontrolü: " + auth + ".\n"
        ws_content += "- Abonelikler: " + (", ".join("`" + s + "`" for s in entry["subscriptions"]) or "yok") + ".\n"
        ws_content += "- İstemci komutları: " + (", ".join("`" + s + "`" for s in entry["commands"]) or "handler içinde komut okuma döngüsü yok") + ".\n"
        if entry["direct_messages"]:
            ws_content += "- Doğrudan yanıtlar: " + ", ".join("`" + s + "`" for s in entry["direct_messages"]) + ".\n"
    write("WEBSOCKET.md", ws_content)

    schema = []
    for path, source in sources.items():
        if not (path.startswith("model/") or path in ("game/score.go", "game/score_summary.go", "game/ranking_fields.go", "game/hub.go", "game/foul.go")):
            continue
        for match in re.finditer(r"^type (\w+) struct \{", source, re.M):
            start = match.end() - 1
            end = balanced_end(source, start)
            fields = []
            for line in source[start + 1:end - 1].splitlines():
                line = line.strip()
                if not line or line.startswith("//"):
                    continue
                line = line.split("//")[0].strip()
                field_match = re.match(r"(\w+)\s+(.+)", line)
                fields.append({"name": field_match[1], "type": field_match[2]} if field_match else {"name": "(gömülü)", "type": line})
            schema.append({"name": match[1], "file": path,
                           "line": source.count("\n", 0, match.start()) + 1, "fields": fields})
    schema_content = "# Veri yapıları ve alan sözlüğü\n\n"
    schema_content += "`model/` yapıları ile temel oyun yapılarının kaynak kodundan çıkarılmış alanlarıdır. "
    schema_content += "Her yapı bir Bolt tablosu değildir; gerçek 13 bucket ana haritada listelenir. "
    schema_content += "Tipler Go gösterimiyle korunur. `db:\"id,manual\"` uygulamanın verdiği anahtarı; `db:\"id\"` otomatik artan anahtarı gösterir. "
    schema_content += "Canlı veri, parola veya API anahtarı değerleri bu dosyaya okunmaz.\n"
    for struct in schema:
        schema_content += "\n## " + struct["name"] + "\n\n" + link(struct["file"], struct["line"]) + "\n\n"
        schema_content += table(["Alan", "Tip / etiket"], [[code(f["name"]), code(f["type"])] for f in struct["fields"]])
    write("VERI-SOZLUGU.md", schema_content)

    packages = []
    for pkg in sorted(set(str(Path(p).parent).replace("\\", "/") for p in gofiles)):
        prod = [p for p in production if str(Path(p).parent).replace("\\", "/") == pkg]
        tests = [p for p in gofiles if p.endswith("_test.go") and str(Path(p).parent).replace("\\", "/") == pkg]
        imports = sorted(set(re.findall(r'"github.com/Team254/cheesy-arena/([^"\n]+)"', "\n".join(sources[p] for p in prod))))
        packages.append({"package": pkg, "source_files": len(prod), "source_lines": sum(len(sources[p].splitlines()) for p in prod),
                         "test_files": len(tests), "test_functions": sum(len(re.findall(r"^func Test\w+\(", read(p), re.M)) for p in tests),
                         "internal_imports": imports})
    inventory_content = "# Kaynak ve bağımlılık envanteri\n\n"
    inventory_content += "Satır sayıları yorumları ve boş satırları içerir. Test fonksiyonu sayısı `func Test...` tanımlarıdır; alt testler sayılmaz. "
    inventory_content += "`*_test.go` dosyaları üretim bağımlılıklarına dahil edilmez; `field/test_helpers.go` ise adına rağmen normal derlemeye dahildir.\n\n"
    inventory_content += table(["Paket", "Go dosyası", "Go satırı", "Test dosyası", "Test fonksiyonu", "İç bağımlılıklar"],
        [[p["package"], p["source_files"], p["source_lines"], p["test_files"], p["test_functions"], ", ".join(p["internal_imports"]) or "—"] for p in packages])
    inventory_content += "\n## Dosya sorumlulukları ve fonksiyon indeksi\n\n"
    for path in production:
        file_functions = [f for f in all_functions if f["file"] == path]
        comments = []
        for line in sources[path].splitlines():
            if line.startswith("package "):
                break
            if line.startswith("//") and not any(s in line for s in ("Copyright", "Author:")):
                comment = line[2:].strip()
                if comment:
                    comments.append(comment)
        inventory_content += "### " + link(path) + "\n\n"
        if comments:
            inventory_content += " ".join(comments) + "\n\n"
        inventory_content += (", ".join(link(path, f["line"], ("(" + f["receiver"] + ")." if f["receiver"] else "") + f["name"]) for f in file_functions) or "Fonksiyon yok; tip, sabit veya üretilmiş enum tanımı.") + "\n\n"
    inventory_content += "## Handler → şablon → tarayıcı dosyaları\n\n"
    inventory_content += "Tablo yalnızca handler gövdesinde doğrudan geçen şablonları listeler. "
    inventory_content += "`renderSettings`, `renderTeams` gibi yardımcıların şablonları sonraki tam şablon envanterinde görünür.\n\n"
    inventory_content += table(["Rota", "Doğrudan kullanılan şablon"],
        [["`" + r["method"] + " " + r["path"] + "`", ", ".join(link(p) for p in r["templates"])] for r in routes if r["templates"]])
    templates = sorted(p for p in tracked if p.startswith("templates/"))
    inventory_content += "\n## Tüm şablonlar\n\n"
    inventory_content += table(["Şablon", "Doğrudan JS/CSS referansları"],
        [[link(p), ", ".join(link(asset.removeprefix("/")) for asset in sorted(set(re.findall(r'(/static/[^"?\s]+\.(?:js|css))', read(p)))))] for p in templates])
    inventory_content += "\n## Uygulamaya ait tarayıcı kodu\n\n"
    jsfiles = sorted(p for p in tracked if p.startswith("static/js/") and "/lib/" not in p)
    inventory_content += table(["Dosya", "Satır", "WebSocket uçları"],
        [[link(p), len(read(p).splitlines()), ", ".join("`" + s + "`" for s in sorted(set(re.findall(r'"([^"\n]*websocket)"', read(p)))))] for p in jsfiles])
    inventory_content += "\n## Çalışma zamanı varlıkları\n\n"
    inventory_content += table(["Tür", "Git içindeki dosya sayısı"], [
        ["schedules/*.csv", sum(p.startswith("schedules/") and p.endswith(".csv") for p in tracked)],
        ["templates/*", len(templates)],
        ["static/js (lib hariç)", len(jsfiles)],
        ["static/css (lib ve fonts hariç)", sum(p.startswith("static/css/") and p.endswith(".css") and "/lib/" not in p for p in tracked)],
        ["static/audio/*", sum(p.startswith("static/audio/") for p in tracked)],
    ])
    write("KAYNAK-ENVANTERI.md", inventory_content)
    commit = subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=ROOT, text=True).strip()
    data = {"source_commit": commit, "packages": packages, "routes": routes,
            "websocket": websocket_entries, "notifiers": notifier_entries, "structs": schema,
            "functions": [{k: v for k, v in f.items() if k != "body"} for f in all_functions]}
    write("envanter.json", json.dumps(data, ensure_ascii=False, indent=2))
    print(json.dumps({"routes": len(routes), "websockets": len(websocket_entries),
                      "notifiers": len(notifier_entries), "structs": len(schema),
                      "functions": len(all_functions), "packages": packages}, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
