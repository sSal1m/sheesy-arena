# Doğrulama sonuçları ve teknik bulgular

Tarih: **29 Ekim 2026**. Kaynak sürümü: **a525aa2**. Ortam: Windows/amd64, Go 1.26.7. İnceleme amacı sistemi haritalamaktır. Aşağıdaki maddeler kaynak gözlemleri, çalıştırılmış kontroller ve önerileri birbirinden ayırır.

## 1. Çalıştırılmış kontroller

| Kontrol | Sonuç | Kapsam |
| --- | --- | --- |
| `go build -o cheesy-arena.exe .` | İlk kurulumda geçti | Windows binary üretimi |
| Yerel ana sayfa, ayarlar, takımlar, match play | İlk kurulumda HTTP 200 | HTML üretimi ve erişim |
| Ana sayfanın yerel JS/CSS/icon bağımlılıkları | İlk kurulumda HTTP 200 | Statik dosya erişimi |
| `/api/arena/websocket` | İlk kurulumda bağlantı açıldı | WebSocket upgrade |
| `go test ./...` | **Başarısız** | `field` ve `web` paketlerinde aşağıdaki testler başarısız; diğer 9 test paketi geçti |
| `go test ./field -run '^TestPlcEStopAStop$' -count=1` | Geçti | İlk PLC test başarısızlığı tek başına tekrarda görülmedi |
| `go test ./web -run '^(TestAllianceStationDisplayWebsocket\|TestAnnouncerDisplayWebsocket\|TestScoringPanelWebsocket)$' -count=1` | **Başarısız** | Alliance Station ve Announcer testleri yine başarısız; Scoring Panel tekrarda geçti |
| `node --test static/js/stream_displays_test.js` | **7/7 geçti** | Twitch/YouTube caption ve boş YouTube yapılandırma senaryoları |
| Haritalama sonu canlı HTTP kontrolü | **9/9 HTTP 200** | Ana sayfa, ayarlar, takımlar, maç kontrolü, ekran yönetimi, qualification JSON, rankings, alliances ve takım PDF raporu |
| Haritalama sonu WebSocket snapshot kontrolü | **Geçti** | `/api/arena/websocket` bağlantısında `matchTiming`, `matchLoad`, `matchTime` alındı; yalnızca mesaj tipleri raporlandı |
| Doküman ve envanter kontrolü | **Geçti** | 1.306 yerel bağlantı, 664 fonksiyon satırı, envanter sayıları, UTF-8, kod blokları ve üretici Python sözdizimi doğrulandı |
| Tekrarlanabilir envanter üretimi | **Geçti** | Aynı kaynak sürümünde iki üretimin SHA-256 çıktıları aynı |

Hedefli Go test komutundaki üç isim gerçek çalıştırmada regex `|` ile ayrılmıştır; tabloda Markdown sütun ayıracını korumak için kaçış gösterimi vardır.

### Tam Go test turunda görülen başarısızlıklar

| Test | Kaynak | Gözlenen durum |
| --- | --- | --- |
| `TestPlcEStopAStop` | [arena_test.go](../../field/arena_test.go#L783) | Stop/Enabled alanları için bazı bool beklentileri uyuşmadı |
| `TestAllianceStationDisplayWebsocket` | [alliance_station_display_test.go](../../web/alliance_station_display_test.go#L65) | Beklenen olaylar okunan sabit mesaj grubunda bulunamadı |
| `TestAnnouncerDisplayWebsocket` | [announcer_display_test.go](../../web/announcer_display_test.go#L120) | Beklenen olay türleri ile okunan olayların sırası/konumu uyuşmadı |
| `TestScoringPanelWebsocket` | [scoring_panel_test.go](../../web/scoring_panel_test.go#L40) | Bağlantının hemen ardından panel kaydı 1 yerine 0 okundu |

Tekrar koşusunda `TestAllianceStationDisplayWebsocket` için `realtimeScore` yerine `arenaStatus`, `TestAnnouncerDisplayWebsocket` için `scorePosted` yerine `realtimeScore` okunması da görüldü.

**Yorum:** Notifier'lar ayrı kanallar, ayrı goroutine'ler ve `reflect.Select` ile iletilir. Testlerin sabit sayıda mesaj okuması veya bağlantıdan hemen sonra yan etkiyi kontrol etmesi zamanlamaya duyarlı olabilir. PLC testi de `time.Now` ve faz eşikleriyle çalışır. Bu, makul bir inceleme yönüdür; başarısızlıkların kesin kök nedeni bu çalışmada ispatlanmamıştır. Tüm hataların yalnızca Windows veya yalnızca test sorunu olduğu iddia edilemez.

İzlenecek iş: testlerde ilgili olayı sınırlı süre boyunca türüne göre bekleme, panel kaydını eventual condition olarak doğrulama ve sahte saat kullanımı değerlendirilmeli; üretim tarafındaki yanlış sıralama/eksik olay olasılığı ayrıca incelenmelidir. Haritalama kapsamında test veya uygulama kodu değiştirilmedi.

## 2. Kaynakta doğrulanan veri ve erişim sınırları

Bu bölüm bir penetrasyon testi değildir. Gerçek WPA, oturum veya servis anahtarı okunmadı; gerçek donanım ya da yayın API'sine komut gönderilmedi.

### B01 — Public mesajlarda tam Team nesnesi

**Kaynakta görülen:** [Team](../../model/team.go) içinde `WpaKey` ve `FtaNotes` exported alanlardır; JSON'dan çıkaran `json:"-"` etiketi yoktur. [GenerateMatchLoadMessage](../../field/arena_notifiers.go#L149), `Teams map[string]*model.Team` ve saha dışı takım listelerini doğrudan taşır. [arenaWebsocketApiHandler](../../web/api.go#L208) admin kontrolü olmadan `MatchLoadNotifier` aboneliği sunar. Salt okunur gösterim uçlarında da aynı tipler kullanılabilir.

**Anlamı:** Yapılandırılmış ve maça atanmış takımlar varsa, bu alanlar API erişimi olan istemcilere serileştirilebilir. Ekranda alanın gösterilmemesi, ağ mesajından çıkartıldığı anlamına gelmez.

**Önerilen öncelik: yüksek.** Public ekran/API için yalnızca gerekli alanları içeren DTO oluşturulmalı; WPA ve özel FTA notları yetkili veri yolunda tutulmalı. Ham model pointer'ının tüm ekranlara gönderilmesi gözden geçirilmeli. Bu çalışma mevcut etkinlikte gerçek değer sızıntısı olduğunu iddia etmez; kodda buna açık veri yolu bulunduğunu saptar.

### B02 — Varsayılan açık yönetim ve tek hesap modeli

**Kaynakta görülen:** [userIsAdmin](../../web/login.go#L63), parola boşsa erişime izin verir. Varsayılan EventSettings parolası boştur. Yetkili kullanıcı adı `admin`'dir. Hakem, scorer ve FTA ayrı kullanıcı türleri değildir. [ServeWebInterface](../../web/web.go#L89), `:8080` üzerinde dinler.

**Anlamı:** `-dev` ile başlatmak web erişimini localhost ile sınırlandırmaz. LAN üzerinden erişim, ağ/firewall koşullarına ve uygulama parolasına bağlıdır. Bu incelemede Windows firewall üzerinde yeni kural açılmadı.

**Önerilen öncelik: yüksek, ağ üzerinden kullanım öncesinde.** İlk yönetici parolasının kurulması ve iş bazlı yetkilerin tanımlanması değerlendirilmeli. Ortamın dış dünyaya açık olduğu varsayılmamalı; mevcut kodun sınırı bu şekilde belgelenmeli.

### B03 — Parola, cookie ve oturum yaşam döngüsü

**Kaynakta görülen:** [login.go](../../web/login.go) parolayı ayardaki değerle doğrudan karşılaştırır. UUID token bbolt'ta tutulur. Cookie oluşturulurken `HttpOnly`, `Secure`, `SameSite` veya açık Path belirlenmez. [GetUserSessionByToken](../../model/user_session.go) `CreatedAt` için son kullanma kontrolü yapmaz. Parola değişiminde oturumlar temizlenir.

**Anlamı:** Kullanıcı/oturum yönetimi basittir; parola hash'i ve süreli oturum politikası bulunmaz. `EventSettings` kaydı ağ ve servis sırlarını da taşır; DB yedeği aynı sırları içerebilir.

**Önerilen öncelik: ağ kullanımında yüksek.** Parola hash'i, cookie nitelikleri, süre aşımı ve oturum iptal akışı tasarlanmalı. Kimlik doğrulama eklendiğinde public HTTP ve WebSocket yollarının da kapsandığı test edilmeli.

### B04 — Şablon kaçış sınırı

**Kaynakta görülen:** [web.go](../../web/web.go), HTML de üreten ortak parse katmanında `text/template` kullanır. Bazı tarayıcı kodları dinamik değerleri `.html(...)` ile işler; örnek [display_shared.js](../../static/js/display_shared.js).

**Anlamı:** `html/template` bağlamına özgü otomatik kaçış sağlanmaz. Kullanıcıdan veya dış servisten gelen metnin hangi HTML/JS bağlamına yerleştirildiği önemlidir. Tek başına bu gözlem, her alanda çalıştırılabilir bir XSS örneği bulunduğunu ispatlamaz.

**Önerilen öncelik: yüksek.** HTML ve CSV/SVG üretimini ayrı değerlendirip HTML için uygun template/escaping kullanımı incelenmeli; kasıtlı HTML alanları açıkça sınırlandırılmalı.

### B05 — Yetkilendirme, HTTP yöntemi ve CSRF

**Kaynakta görülen:** Admin kontrolleri handler'larda dağınıktır. Bazı değişiklik yolları GET olarak kayıtlıdır: takım refresh, WPA üretimi ve bazı yayın tetikleri. [HTTP rota envanteri](HTTP-ROTALARI.md) bunları yöntemleriyle listeler. Formlar için ortak bir CSRF token katmanı görülmez. Gorilla WebSocket upgrader özel CheckOrigin tanımlamaz; kütüphanenin varsayılan origin davranışı geçerlidir.

**Anlamı:** Menünün bir bağlantıyı gizlemesi sunucu yetkisi değildir. HTTP CSRF önlemleri ile WebSocket origin kontrolü ayrı konulardır; ikisini de tamamen yok saymak veya aynı önlem kabul etmek yanlış olur.

**Önerilen öncelik: orta/yüksek.** Değişiklik yapan GET uçları, CSRF, FTA koşulları ve rol gereksinimleri merkezi erişim tasarımıyla gözden geçirilmeli.

## 3. Dayanıklılık ve mimari sınırlar

| Bulgular | Kaynakta görülen | Etki ve önerilen inceleme |
| --- | --- | --- |
| B06 — Disk senkronizasyonu | [database.go](../../model/database.go): `NoSync:true` | Ani güç kesintisi dayanıklılığı, performans ve yedek stratejisi birlikte değerlendirilmeli |
| B07 — Çok adımlı sonuç işlemi | [commitMatchScore](../../web/match_play.go#L402): farklı model metotlarıyla art arda kayıt | Ara hatada sonuç/maç/sıralama tutarsızlığı olabilir; atomik kayıt veya yeniden hesaplama/kurtarma planı incelenmeli |
| B08 — Mesaj düşmesi | [notifier.go](../../websocket/notifier.go): 5 kapasiteli kanal, non-blocking send | Yavaş istemci olay kaçırabilir; olay ile snapshot ayrımı ve yeniden senkronizasyon ihtiyacı |
| B09 — Paylaşılan Arena alanları | HTTP/WS, DS ve PLC goroutine'leri aynı nesnelere erişir | Bazı alt sistemlerde mutex var; tüm Arena için genel senkronizasyon yok. Yarış riski analiz edilmeli; bu çalışmada race detector koşulmadı |
| B10 — Bellek içi operasyon durumu | `NewArena` boş test maçı başlatır | Yeniden başlatma, aktif maçın kaldığı yerden sürmesi değildir; operasyonel kurtarma akışı gerekir |
| B11 — Dış servis timeout farkları | TBA'da açık HTTP timeout yok; Nexus 30 s, AP 3 s | Yükleme veya yayın hatasının etkisi değişir; süre sınırı, iptal ve yeniden deneme politikaları incelenmeli |
| B12 — Cihaz taşıma güveni | Ana switch Telnet; SCC'de `InsecureIgnoreHostKey` | Güvenilir saha ağı varsayımı var; cihaz/protokol değişiminde host doğrulama ve bağlantı koruması değerlendirilmeli |
| B13 — Dosya konumuna bağlı dağıtım | Static/templates/schedules/font dosya sisteminden okunur | Exe tek başına taşınamaz; dağıtım paketi ve çalışma dizini korunmalı |
| B14 — Restore işlem sınırı | [restoreDbHandler](../../web/setup_settings.go): DB kapanır/değişir, ayarlar yüklenir | Canlı Arena nesnesi tümüyle yeniden kurulmaz; concurrent istekler ve aktif maç durumu için operasyon koruması incelenmeli |
| B15 — HTTP sunucu yaşam döngüsü | [web.go](../../web/web.go): `http.ListenAndServe`; [main.go](../../main.go): Arena sonsuz döngü | Özel server timeout'ları, context tabanlı kapatma ve graceful shutdown akışı görünmüyor |
| B16 — Paketler arası yoğun bağımlılık | Web ve Arena iş kurallarını birlikte yönetir | Yeni sezon/rol/API çalışmalarında değişiklik alanı geniş; sonuç servisi ve public DTO ayrımı yararlı olabilir |

SCC/Telnet ve TBA imza davranışları mevcut protokollerin tespitidir; bu belgede alternatif cihaz desteğinin hazır olduğu veya yapılmış olduğu iddia edilmez. Kaynak içindeki sabit servis anahtarı değerleri belgeye kopyalanmadı.

## 4. Değişiklik yapılmadan açık bırakılan konular

- Gerçek robotların DS protokol uyumu ve etkinleştirme/durdurma akışı sahada denenmedi.
- Fiziksel PLC programı, sensör kablolaması, ArmorBlock'lar ve motor çıkışları doğrulanmadı.
- VH-113 AP, ana switch, SCC, LED ve Cypress tabelalarıyla bağlantı denenmedi.
- TBA/Nexus yayınlama ve HyperDeck/Companion komutları gerçek cihaz/etkinlik üzerinde çalıştırılmadı.
- Tarayıcı ekranlarının tüm çözünürlüklerde görsel incelemesi yapılmadı; bu çalışma arayüz tasarım denetimi değildir.
- Test başarısızlıklarının kesin kök nedeni çözülmedi; uygulama davranışında onarım yapılmadı.
- `go test -race`, güvenlik tarayıcısı, yük testi ve uzun süreli saha dayanıklılık testi çalıştırılmadı.

## 5. Önerilen çalışma sırası

1. Public veri mesajlarını güvenli DTO'lara ayır; parolalı ağ kullanımı sınırını netleştir.
2. Tekrarlanan WebSocket test başarısızlıklarını olay türüne göre bekleme/senkronizasyon açısından teşhis et; PLC testindeki zaman etkisini ayır.
3. Sonuç kaydetme, yedek ve restore için tutarlılık/kurtarma senaryolarını tanımla.
4. Rol/yetki, cookie, CSRF ve HTML kaçışını birlikte ele al.
5. Gerçek saha için DS/PLC/ağ/LED/tabela kabul testlerini cihazlarla uygula.
6. Yeni ürün veya sezon özelleştirmesine, bu sınırlar ve etki haritası üzerinden başla.
