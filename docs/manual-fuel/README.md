# Manual FUEL Control — kullanım ve teknik teslim

Ek veri kaynağı: [Dual Arduino Ethernet entegrasyonu](../arduino/README.md). Otomatik kaynağın PLC veya Arduino olduğu ittifak bazında belirlenir; aşağıdaki manuel düzeltme ve kurtarma kuralları korunur.

Bu özellik 2026 REBUILT için mevcut Cheesy Arena skor hattına entegredir. Paneldeki değerler **puanlanan FUEL** sayılarıdır. Etkin değerler mevcut `game.Hub.ShiftCounts` alanına yansıtılır; `Score.Summarize`, bonus/RP, canlı ekranlar, MatchResult, sıralama, playoff ve TBA/Nexus sonuç akışları kullanılmaya devam eder.

## Erişim ve Windows'ta çalıştırma

Panel: **http://localhost:8080/panels/manual_fuel**. Başka cihazdan `http://BILGISAYARIN-YEREL-IP-ADRESI:8080/panels/manual_fuel` kullanılır. Bağlantı **Panel → Manual FUEL Control** menüsünde ve maç kontrol sayfasındadır. Mevcut yönetici parolası tanımlıysa admin oturumu gerekir. Parola boşsa mevcut yönetim ekranlarıyla aynı açık erişim uygulanır. Yeni hesap veya rol sistemi eklenmedi.

Klasörde PowerShell açın:

```powershell
.\Baslat.ps1
# Durdurmak için:
.\Durdur.ps1
```

`Baslat.cmd` / `Durdur.cmd` dosyaları da kullanılabilir. Başlatıcı `-dev` ile çalışır. Bu seçenek mevcut geliştirme ağı davranışıdır; gerçek robot/saha kurulumu için kök README'deki ağ ve PLC ayarları uygulanmalı ve program `-dev` olmadan başlatılmalıdır. Manual Mode robot bağlantısı, bypass veya maç başlatma koşullarını kaldırmaz.

Kaynak koddan derleme (Go sürümü `go.mod`: 1.26.7; Node yalnız JavaScript testleri için):

```powershell
.\Durdur.ps1
go build -o cheesy-arena.exe .
.\Baslat.ps1
```

Sunucuyu repo kökünden çalıştırın; `templates/`, `static/` ve `event.db` göreli yollarla kullanılır. Etkinlik verileri ve işlem günlüğü aynı `event.db` içindedir. Uygulama kapalıyken dosyayı kopyalayabilir veya çalışan uygulamanın mevcut veritabanı yedekleme ekranını kullanabilirsiniz.

## Operatör iş akışı

1. Maçı normal maç kontrol ekranından yükleyin. Sayımlar sıfır ve maç başlamamışken modu seçin. Seçilen mod iki ittifaka ortaktır ve sonraki maçlarda korunur.
2. **Manual**: yalnız panelin manuel katmanı puanlanır. PLC olmadan çalışır. Diğer PLC güvenlik, ışık ve motor işlevleri korunur.
3. **PLC + Manual Correction** (varsayılan): PLC'nin ham sayımı + manuel düzeltme kullanılır. Örneğin ham 25 ve düzeltme +3, etkin 28 üretir; ham 27 olduğunda etkin 30 olur. Düzeltme PLC fark hesabına girmez.
4. İsteğe bağlı operatör etiketi girin. Her karttaki ±1 / ±5 sunucuya delta gönderir. Onay gelene kadar toplam yerel olarak artırılmaz. Ham sayım ve düzeltme ayrıca gösterilir.
5. AUTO ve pause sırasında AUTO; TELEOP ve PostMatch sırasında iki faz da düzenlenebilir. Önceki AUTO düzeltmesi, TELEOP'ta belirlenmiş auto kazananını/Hub sırasını/DS oyun verisini değiştirmez. Sayı sıfırın altına inecekse işlem bütünüyle reddedilir.
6. Sonucu mevcut **Commit and Post** akışından onaylayın. PLC modunda etkin PLC varsa maç sonundaki mevcut 3 saniyelik sayım grace period'u beklenir. Kaydedilen sonuç bu panelden değiştirilemez; mevcut Match Review kullanılır.

Her yeni maç yüklemesi, aynı maçın tekrar oynanması ve kurtarma için yeni bir oynanış bağlamı oluşturulur. Önceki bağlamdan geciken yeni komutlar reddedilir. Başarıyla kaydedilmiş eski bir isteğin tekrarı yalnız önceki kabulü uzlaştırır, yeniden skor uygulamaz. Mod canlı maçta değiştirilemez.

Bağlantı kesildiğinde kontroller kapanır. Onayı belirsiz komutlar sekmenin `sessionStorage` kuyruğunda aynı RequestId ile tutulur; bağlantı/yenileme sonrasında sunucuyla uzlaştırılır. Yeni bağlam görülürse eski kuyruk yeniden uygulanmaz. Kuyruk en fazla 32 işlemdir. Eski revizyonlar ekranı geriye götürmez; bağlantı yeniden kurulurken sunucunun yeni başlangıç snapshot'ı esas alınır.

## PLC kesintisi ve yeniden baz alma

Başarılı register okuması; sayaçlar, okuma zamanı, sıra, bağlantı nesli ve reset nesliyle tek snapshot olarak yayımlanır. Yeni maç resetinden sonra taze sıfır görülmeden sayaç kabul edilmez. Reset öncesinde başlamış bir register okumasının reset sonrasında dönmesi de epoch kontrolüyle engellenir.

500 ms'den eski/sağlıksız örnek, sayaç gerilemesi veya bağlantı/reset nesli değişmesi otomatik FUEL alımını askıya alır. Etkin skor ve manuel düzeltmeler korunur. Ekran nedeni gösterir; bağlantının sağlıklı hale gelmesi kendiliğinden puan eklemez. Yetkili operatör askıdaki ilgili ittifakı **Rebase PLC** ile taze sağlıklı örneğe yeniden baz alabilir. Bu işlem puan eklemez; kaçan FUEL panelden düzeltilebilir. Çalışan, askıda olmayan sayaca rebase reddedilir.

**Donanım sınırı:** mevcut PLC protokolünde maç kimliği yoktur. Eski bir değer donanım tarafından yeni, sağlıklı ve artan bir okuma gibi sunulursa bunu her durumda ayırt etmek mümkün değildir. Yeni maçta sıfır kapısı, nesil kontrolü ve belirsizlikte askıya alma bu riski sınırlar; fiziksel saha kabul testinin yerini almaz.

## Kalıcılık, hata ve kurtarma

Manuel işlem, RequestId tekrarlama kaydı ve skor checkpoint'i aynı bbolt transaction'ında yazılır. DB'nin mevcut `NoSync` ayarına rağmen kabulden önce açıkça `Sync` çağrılır. Yazım hatasında RAM'deki manuel değer değiştirilmez ve başarı gönderilmez. Yazımın kalıcı olup olmadığı belirsizse tekrar istekleri dahil değişiklikler kapanır; sunucu yeniden açılıp kalıcı kayıt uzlaştırılmalıdır.

PLC ve diğer skor değişiklikleri 100 ms aralıklı, tek sıralı arka plan yazıcısıyla checkpoint edilir. Snapshot işlem kilidini aldıktan sonra yakalanır; eski bir arka plan kaydı yeni manuel düzeltmeyi silemez. Disk gecikmesi Arena'nın skor kilidi altında bekletilmez. Ani kesintide **son kalıcı PLC/diğer skor checkpoint'inden sonraki veriler kaybolabilir**; süre disk ve sistem yüküne de bağlıdır. Panel son başarılı checkpoint zamanını gösterir. Başarısı bildirilmiş manuel işlemler senkron kalıcı kayıt kapsamındadır.

Tamamlanmamış oynanış yeniden açılışta **PostMatch inceleme** olarak kurtarılır. Robotlar, saat, motorlar otomatik devam etmez; bypass'lar açılmaz. Operatör düzeltip sonucu onaylayabilir veya mevcut discard akışıyla iptal edebilir. Tarayıcı bağlamı yenilenir; sonuç kimliği aynı oynanışa bağlı kalır. Kısmi sonuç kaydı sonrasında yeniden onay, aynı MatchResult/PlayNumber'ı kullanır. Önceden kaydedilmiş eski sonuçlarda FUEL metadata zorunlu değildir.

Aktif tam skor düzenleyicisi bağlam ve revizyon kontrolüyle eşzamanlı değişikliği ezmeyi engeller. Match Review, sunucuya ait ham sayım metadata'sını korur ve etkin skor değiştiğinde düzeltme katmanını uzlaştırır. Veritabanını silme/geri yükleme, bekleyen maç varken reddedilir; bu işlemler checkpoint yazıcısıyla aynı bakım bariyerini kullanır.

## Teknik harita ve değişen dosyalar

| Alan | Dosyalar ve sorumluluk |
| --- | --- |
| Skor koordinatörü | `field/manual_fuel.go`: ham/manüel katman, faz dağılımı, delta doğrulama, idempotence, PLC kabulü, snapshot, kalıcılık ve kurtarma |
| Diğer skor yazarları | `field/score_edit.go`, `field/fuel_maintenance.go`: hakem/skor editörü, sonuç değiştirme ve DB bakım bariyeri |
| Maç döngüsü ve yayın | `field/arena.go`, `arena_notifiers.go`, `arena_leds.go`, `team_sign.go`: yaşam döngüsü, mevcut hesaplamalara projeksiyon, kilitli bağımsız snapshot'lar |
| Veri modeli | `model/fuel.go`, `model/match_result.go`: ham/düzeltme metadata, checkpoint/işlem bbolt bucket'ları, sonuç deep clone |
| PLC | `plc/fuel.go`, `plc/plc.go`: tutarlı snapshot ve reset/bağlantı nesilleri; mevcut I/O korunur |
| HTTP/WS | `web/manual_fuel.go`, `web/web.go`: panel, WebSocket ve geçmiş endpoint'i; her komutta mevcut admin yetkisi |
| Mevcut iş akışları | `web/match_play.go`, `match_review.go`, `referee_panel.go`, `scoring_panel.go`, `setup_settings.go`: ortak skor yazma sınırı, sonuç snapshot/onay ve DB bakım entegrasyonu |
| Arayüz | `templates/manual_fuel.html`, `static/css/manual_fuel.css`, `static/js/manual_fuel*.js`; `templates/base.html`, `match_play.html` bağlantıları |
| Testler | `field/manual_fuel_test.go`, `model/fuel_test.go`, `plc/fuel_test.go`, `web/manual_fuel_test.go`, `partner/manual_fuel_test.go`, `static/js/manual_fuel_client_test.js` ve mevcut fake/olay sıralama testleri |

Temel yol:

```text
PLC snapshot → doğrulama/cursor → ham Hub kovaları ─┐
                                                 ├→ etkin Hub.ShiftCounts
admin WS delta → kalıcı transaction → düzeltmeler ─┘
  → mevcut Score.Summarize / RealtimeScoreNotifier
  → ekranlar / sonuç snapshot → MatchResult → sıralama / playoff / TBA / Nexus
```

Manuel AUTO auto kovasına; TELEOP artırımı mevcut puanlanan vardiyaya, Hub pasifken son başlamış aktif TELEOP vardiyasına, maç sonrasında endgame kovasına yazılır. Azaltım ilgili fazın puanlanan kovalarından yeniden eskiye uygulanır. Böylece vardiya dökümü ile etkin TELEOP toplamı aynı kaynaktan gelir. Yeni bir puanlama motoru veya bağımlılık eklenmedi.

WebSocket `/panels/manual_fuel/websocket`; komutlar `adjustFuel`, `setFuelMode`, `rebasePlcFuel`, `getFuelState`; yanıtlar `manualFuelState`, RequestId içeren `fuelResult`. Delta komutu RequestId, MatchId, AttemptId, Alliance, Phase ve Delta taşır. Aynı kimlik/farklı içerik reddedilir. Audit endpoint'i `/panels/manual_fuel/history?attemptId=...`; panel son 100 işlemi gösterir, DB günlüğü tüm kabul edilen işlemleri saklar. Ortak admin hesabı gerçek kişi kimliğini belirleyemez; Actor oturumu ve isteğe bağlı Operator etiketi ayrı kaydedilir. Günlük zaman, maç/bağlam, faz, ittifak, delta, önce/sonra, mod ve PLC baz/sıra bilgilerini içerir.

## Otomatik doğrulama — 29 Ekim 2026

```powershell
go test ./...
node --test static/js/*test.js
go build -o cheesy-arena-next.exe .
go vet ./...
```

- Tam Go paketi testleri geçti; Node 17/17 geçti; Windows executable derlendi.
- Dört sayaç ±, negatif/overflow reddi, PLC'siz manuel, 25+3→28 ve 27+3→30, ham cursor bağımsızlığı, bonus/RP/G206 ve TBA vardiya/toplam tutarlılığı otomatik test edildi.
- AUTO→TELEOP/geç AUTO/frozen winner, pasif Hub dağılımı, maç sonu grace period, sonuç onayı/kilit, yeni/aynı maç bağlamı, restart/kurtarma ve yarım onay senaryoları test edildi.
- Eşzamanlı operatörler, tekrarlanan/değiştirilmiş istek, eski bağlam, WS yeniden bağlanma, ters revizyon, outbox/yenileme, DB hatası ve kalıcı yazım sonrası ACK öncesi kesinti test edildi.
- PLC reset/reconnect/gerileme/eski/sağlıksız örnek, ilk sıfır kapısı, rebase, reset öncesi uçuşta okuma, mevcut güvenlik I/O, eski sonuç uyumluluğu ve admin oturum iptali test edildi.
- Başlangıç referansında Windows'ta görülen WebSocket olay sırası/kayıt zamanlaması ve PLC testinin tam sıfır zamanına denk gelmesi kaynaklı kararsızlıklar giderildi. Testler kapatılmadı; olayların tamamını bekleme ve sınırdan uzak zaman kullanıldı. İlgili testler ayrıca beş kez tekrarlandı ve geçti.
- `go vet ./...` temiz değil: mevcut testlerde unkeyed literals, mevcut IPv4 `net.Dial` adres oluşturma ve `model/table_test.go` içindeki eski struct tag uyarıları var. Bu ilgisiz alanlar değiştirilmedi. `go vet -composites=false -hostport=false -structtag=false ./...` geçti; bu komut tam vet'in yerine başarılı gibi sunulmamalı.
- **Race detector doğrulanamadı:** SHA256 doğrulanmış resmi Zig 0.14.1 ile CGO/race derlemesi yapıldı; Windows senkronizasyon import kütüphanesi eklendiğinde executable bağlandı, fakat ThreadSanitizer testten önce `failed to allocate ... (error code: 87)` ile kapandı. Eşzamanlılık testleri normal Go çalışmasında geçti; race temizliği iddia edilmiyor. Güvenlik/ASLR ayarları değiştirilmedi.

Uygun CGO derleyicisi ve çalışan ThreadSanitizer ortamında ayrıca çalıştırın:

```powershell
$env:CGO_ENABLED = '1'
# CC, Windows Go race desteği olan mevcut C derleyicinize ayarlanmalıdır.
go test -race ./field ./plc ./model ./web ./partner
```

Testler simüle PLC/WS ve ayrı test veritabanlarıyla yapıldı. Fiziksel PLC, gerçek sensör/Hub, robot DS, ağ kesintisi, motor/ışık, Nexus/TBA canlı servis kabulü bu bilgisayarda denenmedi. Sahaya çıkmadan reset/sıfır kapısı, kablo kesip takma/rebase, faz sınırları/grace, E-stop/A-stop ve sonuç yayını gerçek donanımla doğrulanmalıdır.

## Ekran görüntüleri

İzole test sunucusunda örnek skorlarla alındı; gerçek etkinlik veritabanına örnek skor yazılmadı.

- [Masaüstü, 1440 px](desktop.png)
- [Tablet, 820 px](tablet.png)
