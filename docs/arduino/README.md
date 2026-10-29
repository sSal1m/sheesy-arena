# Dual Arduino Ethernet FUEL — kurulum ve teknik teslim

Arduino olayları mevcut Manual FUEL Control sisteminin **ham otomatik sayım** katmanına entegredir. Manuel düzeltmeler ayrı kalır; etkin `Hub.ShiftCounts`, `Score.Summarize`, bonus/RP, canlı WebSocket ekranları, MatchResult, sıralama ve mevcut TBA/Nexus sonuç hattı kullanılır. Ayrı bir puanlama motoru veya maç zamanlayıcısı eklenmedi.

## Erişim

- Yönetim: **http://localhost:8080/panels/arduino** — Panel → Arduino FUEL Devices.
- Manuel düzeltme: **http://localhost:8080/panels/manual_fuel**.
- Arduino bağlantıları: bilgisayardaki **TCP 1760**. Bu port mevcut DS 1160/1750, HTTP 8080 ve PLC Modbus bağlantısından ayrıdır.
- Yönetim sayfası, durum JSON'u ve simülasyon mevcut admin yetkilendirmesini kullanır. Parola boşsa mevcut açık yönetim erişimi korunur. Cihazlar admin parolası kullanmaz; birbirinden farklı cihaz anahtarları kullanır.

## Kaynak sahipliği ve oyun davranışı

Her ittifakın tek otomatik kaynağı vardır: `plc` veya `arduino`. İkisi aynı hedefi izlerken toplanmaz ve otomatik failover yapılmaz. Örneğin RED=`arduino`, BLUE=`plc` mümkündür. Arduino seçili ittifakta PLC'nin FUEL girdisi atlanır; E-stop/A-stop, motor, ışık ve diğer PLC işlevleri korunur. PLC seçili ittifakta Arduino olayları skor değiştirmez. Manual Mode her iki otomatik kaynağı da FUEL açısından kapatır.

Kaynaklar `arduino.local.json` üzerinden, sunucu açılırken seçilir; canlı maçta değiştirilemez. Sonraki her maçın metadata'sına kaynak sahipliği kopyalanır. Kurtarılan bekleyen maç kendi eski sahipliğini korur; yapılandırma yeni maça uygulanır. Panel **current / configured source** ayrımını gösterir. Eski DB/sonuçlarda alan yoksa kaynak PLC kabul edilir. Arayüzde eski `plc_correction` modunun adı artık **Automatic + Manual Correction**; kayıt değeri geriye uyumluluk için aynıdır.

Sunucu, olayın kabul anındaki gerçek Arena durumunu ve mevcut `Hub.ScoringShiftAt` / `GetShiftCount` kurallarını kullanır. Bu yöntem PLC'nin kullandığı aynı AUTO, transition, vardiya, endgame ve grace period sınıflandırmasını paylaşır. Fiziksel HUB aktiflik bilgisi ile grace nedeniyle hâlâ puanlanabilen durum ayrı bayraklardır. Pasif HUB olayı reddedilir; cihazın fiziksel detected sayacında görülebilir. AUTO kazananı belirlenince geç AUTO manuel düzeltmesi mevcut davranış gereği HUB sırasını değiştirmez.

Maç öncesi, Manual Mode, kurtarma, kapanmış/onaylanmış sonuç veya geçersiz fazda otomatik puan eklenmez. Normal maç sonunda ve abort sonrasında mevcut 3 saniyelik grace sınırı uygulanır. Arduino seçiliyse PLC takılı olmasa da sonuç onayı bu pencerenin kapanmasını bekler. Sonuç onayı aynı işlem bariyerini kullanır; sonradan gelen olay skoru değiştiremez.

Olayın kabulü kalıcı yazım tamamlanınca kesinleşir. AUTO kazananı o anda Arena'ya uygulanmış skordan belirlenir; faz sınırındaki disk gecikmesi nedeniyle sonradan onaylanan AUTO olayı skoru düzeltebilir fakat belirlenmiş HUB sırasını ve DS oyun verisini geriye dönük değiştirmez. Ağ/disk sınırında bu ayrım saha kabul testinde ayrıca kontrol edilmelidir.

## Ağ ve switch

| Bileşen | Adres / bağlantı |
| --- | --- |
| Windows saha bilgisayarı | `10.0.100.5/24`, TCP 1760 dinler |
| Cisco switch | `10.0.100.3`, VLAN 100 |
| RED Uno | `10.0.100.201/24`, Gi0/23, `fuel-red-01` |
| BLUE Uno | `10.0.100.202/24`, Gi0/24, `fuel-blue-01` |

Repo ve [upstream referans konfigürasyon](https://github.com/Team254/cheesy-arena/blob/main/switch_config.txt) Gi0/23 ve Gi0/24'ü zaten VLAN 100 access port yapıyor. `switch_config.txt` değiştirilmedi. `ConfigureTeamEthernet` yalnız takım SVI/DHCP VLAN'ları 10–60'ı düzenliyor; fiziksel Gi0/23–24 veya VLAN 100'e yazmıyor. Otomatik test mevcut komutları ve statik port bloklarını kontrol ediyor.

Fiziksel switch üzerinde doğrulama:

```text
show interfaces Gi0/23 switchport
show interfaces Gi0/24 switchport
show vlan brief
show ip interface brief
```

Yalnız portlar yanlışsa uygulanacak dar kapsamlı düzeltme:

```text
configure terminal
interface range GigabitEthernet0/23 - 24
 switchport mode access
 switchport access vlan 100
 no shutdown
end
```

Robot VLAN'larını, DHCP havuzlarını veya switch'in geri kalan ayarlarını değiştirmeyin. İki cihaz aynı VLAN/subnet içinde olduğundan aralarındaki FUEL trafiği için router gerekmez. Bağlı Windows Ethernet adaptörünü belirleyip `10.0.100.5/24` ayarlayın; mevcut başka bir adaptörü varsayımla değiştirmeyin. Bu geliştirme sırasında bilgisayarda `10.0.100.5` atanmış bir IPv4 adaptörü bulunmadı; gerçek NIC ve switch yapılandırması uzaktan değiştirilmedi.

Gerekliyse yönetici PowerShell'de, yalnız saha ağından gelecek bağlantılar için firewall kuralı:

```powershell
New-NetFirewallRule -DisplayName 'Cheesy Arena Arduino FUEL' -Direction Inbound -Action Allow -Protocol TCP -LocalAddress 10.0.100.5 -LocalPort 1760 -RemoteAddress 10.0.100.201,10.0.100.202 -Program "$PWD\cheesy-arena.exe"
```

## Yapılandırma ve anahtar üretimi

Bu bilgisayarda `arduino.local.json` ve RED/BLUE özel `DeviceConfig.h` dosyaları üretildi. **Enabled=false** bırakıldı; fiziksel donanım/ağ doğrulanmadan otomatik kaynağı devreye sokmaz. Dosyalar Git tarafından dışlanır, özel anahtarlar repoya gönderilmez. Genel şablon `arduino.example.json` içinde anahtar bulunmaz.

Yeni kurulumda (mevcut dosyaları tekrar üretmeyin):

```powershell
.\tools\New-ArduinoConfig.ps1
# İki bağımsız atış kanalı varsa üretim sırasında: -SensorChannels 2
```

Script mevcut sunucu anahtar dosyasını ezmez. Her cihaz için 32 bayt kriptografik rastgele anahtar ve farklı MAC üretir; aynı anahtarı sunucu JSON'u ile ilgili firmware başlığına yazar. Dosyaları gizli tutun. Yeniden anahtar üretimi cihazın da tekrar programlanmasını gerektirir.

- Sunucu: `arduino.local.json`; `Enabled`, `Listen`, `Devices[].Id/Alliance/IP/Key/Source`.
- RED: `firmware/local/red/DeviceConfig.h`.
- BLUE: `firmware/local/blue/DeviceConfig.h`.
- Yalnız donanım ve adresler hazır olduğunda JSON'da `Enabled=true` yapın. Kaynaklar `arduino` ise iki Arduino sayımı kullanılır. PLC'yi koruyup yalnız cihaz haberleşmesini denemek için ilgili `Source` değerini `plc` bırakabilirsiniz.
- Başka dosya: `cheesy-arena.exe -arduino-config C:\saha\arduino.json`.
- FUEL portu değiştirilirse JSON `Listen` ile iki firmware başlığındaki `SERVER_PORT` birlikte değiştirilmelidir. Varsayılan 1760'tır; mevcut Arena portları 1160/1750/8080 ve ayrıcalıklı portlar kabul edilmez.
- Dosya yoksa sistem Arduino dinleyicisi açmadan eski PLC/manüel davranışıyla çalışır. Etkin fakat hatalı dosyada veya dinleme adresi açılamadığında program açık hata ile durur; sessizce başka bir kaynağa geçmez.

## Firmware: derleme ve yükleme

Kaynak `firmware/FuelEthernet/FuelEthernet.ino`. [Arduino CLI](https://arduino.github.io/arduino-cli/latest/installation/), Arduino AVR core 1.8.6, [Ethernet 2.0.2](https://github.com/arduino-libraries/Ethernet) ve [Crypto 0.4.0](https://github.com/rweather/arduinolibs/tree/master/libraries/Crypto) ile Uno hedefinde derlendi. Go tarafına bağımlılık eklenmedi; firmware'in iki kütüphanesi ağ ve HMAC için gereklidir.

```powershell
arduino-cli core update-index
arduino-cli core install arduino:avr@1.8.6
arduino-cli lib install Ethernet@2.0.2 Crypto@0.4.0
.\tools\Build-Arduino.ps1
arduino-cli board list
# COM numaralarını bağladığınız gerçek kartlara göre değiştirin:
arduino-cli upload --fqbn arduino:avr:uno --port COM5 --input-dir .\firmware\local\build\red\FuelEthernet\output
arduino-cli upload --fqbn arduino:avr:uno --port COM6 --input-dir .\firmware\local\build\blue\FuelEthernet\output
```

RED firmware'ini BLUE kartına veya tersini yüklemeyin. Sunucu kimlik, IP ve anahtar eşleşmesini zorunlu tutar. Bu bilgisayardaki doğrulama araçları `%TEMP%\sheesy-arduino-tools` altında ayrı veri/kütüphane dizinlerine kuruldu; mevcut Arduino kurulumu değiştirilmedi. Standart kurulumu olan makinede yukarıdaki komutlar yeterlidir.

`Build-Arduino.ps1` kaynak ve ilgili özel başlığı ayrı RED/BLUE sketch dizinlerine kopyalar. Her iki `.hex` ve `.elf` dosyası `firmware/local/build/.../output` altında üretilir. Script yalnız derler; bağlı cihaza kendiliğinden yüklemez.

### Shield ve sensör uyarlaması

Ethernet sürücüsü standart Arduino Ethernet kütüphanesidir; W5100/W5200/W5500 ailesi için tasarlanmıştır. Gerçek shield modeli henüz seçilmediğinden fiziksel uyumluluk doğrulanmadı. ENC28J60 aynı sürücü değildir; bu kart seçilirse ağ adaptörü/kütüphanesi ayrıca uyarlanıp Uno bellek ve iletişim testleri yeniden yapılmalıdır.

Özel başlıkta CS pini, SD CS, MAC/IP, sunucu/port, sensör pinleri, kanal sayısı, aktif polarite, pull-up ve debounce süresi merkezîdir. `sensorEdge` giriş adaptörü, ağ/maç mantığından ayrıdır. Varsayılan **tek kanal D2**, aktif LOW, pull-up, 1000 µs debounce'dur. İki bağımsız top geçiş kanalı varsa D2/D3 ve `SENSOR_CHANNELS=2` kullanılır. Aynı topu gören iki sensör iki bağımsız kanal gibi sayılmamalıdır; böyle bir donanımda giriş adaptörü tek top olayı üretmelidir.

Her aktif kenar bir olay üretir; ışık kesilmesi boyunca tekrar sayılmaz. Debounce, tek başına her sensörün doğru top ayrımını garanti etmez. Pulse genişliği, iki top aralığı ve mekanik düzen gerçek sensörle ölçülmelidir. ISR içinde ağ, disk, `String` veya gecikme yoktur. 16 girişli ring buffer'ın kullanılabilir kapasitesi 15 olaydır; ACK gelene kadar olay tutulur. Ana döngü HMAC ve Ethernet işler; kısmi/birleşik TCP çerçeveleri desteklenir.

Şebeke kopukken, eski durumla veya kuyruk doluyken algılanan olaylar detected/dropped telemetrisiyle raporlanır. Sınırsız offline saklama veya kayıpsız sayım iddiası yoktur. Bağlam değişiminde kuyruktan çıkarılan, ACK'i henüz alınmamış fakat sunucuda kabul edilmiş bir olay da dropped'a dahil olabilir; bu sayaç doğrudan “eksik puan” değildir. Manuel düzeltmeden önce kabul günlüğüyle karşılaştırın. Cihaz yeniden başlatılırsa RAM kuyruğu kaybolur. EEPROM'un ilk 128 baytı, her boot'ta bir kez ilerleyen 16 kayıtlı boot kimliği günlüğüne ayrılmıştır; her top için EEPROM yazılmaz. Reset/anahtar/EEPROM işlemlerini maç sırasında yapmayın.

## Donanım bağlantı şeması

[Ağ ve sensör bağlantı şeması](wiring.svg).

| Bağlantı | Uygulama |
| --- | --- |
| RED Ethernet | switch Gi0/23, VLAN 100 |
| BLUE Ethernet | switch Gi0/24, VLAN 100 |
| Uno ↔ shield | kartın kendi shield/ICSP konnektörü; seçilen modelin pinout'una göre |
| Ethernet CS | varsayılan D10; SD kullanılmaz, D4 HIGH |
| Sensör kanal 1 | uyumlu dijital çıkış → D2; ikinci bağımsız kanal varsa D3 |
| Lojik referans | sensör arayüzünün lojik GND'si → Uno GND, izolasyon gerekiyorsa uygun arayüz üzerinden |

**Sensör beslemesi modele bağlıdır.** Belirsiz sensörü doğrudan Uno pinine bağlamayın. Giriş Uno'nun elektriksel sınırlarına uygun lojik seviyeye çevrilmiş olmalıdır. 12/24 V endüstriyel çıkış için uygun optik izolasyon/level shifter gerekir; besleme gerilimi GPIO'ya verilmez. Pull-up ve aktif LOW varsayımı ancak sensör çıkış tipi uygunsa geçerlidir. D10–D13/ICSP ve D4 shield tarafından kullanılır; sensör için ayrılan D2/D3 ile çakıştırmayın. Gerçek sensör modeli verilmediği için konektör renkleri ve besleme uçları uydurulmadı.

## Protokol ve dayanıklılık

F2 v1, sabit 88 bayt çerçeve kullanır. Her tamsayı big-endian; `io.ReadFull` paket sınırlarından bağımsız okur. Her çerçevenin ilk 72 baytı cihaz anahtarıyla HMAC-SHA256'nın 16 bayt etiketiyle doğrulanır. Bu kimlik/bütünlük doğrulamasıdır; trafik şifrelenmez.

| Bayt | Alan |
| --- | --- |
| 0–3 | `F2`, sürüm 1, mesaj tipi |
| 4–7 | bağlantı içi sıra |
| 8–23 | UUID oynanış bağlamı |
| 24–27 | kısa ömürlü durum lease'i |
| 28–31 | detected; STATE içinde MatchId |
| 32–35 | bayraklar / ACK sonucu / HELLO firmware sürümü |
| 36–39 | olay numarası / boot / STATE faz+vardiya / dropped |
| 40–55 | sıfır dolgulu sabit cihaz kimliği |
| 56–71 | sunucunun rastgele oturum nonce'u |
| 72–87 | HMAC etiketi |

HELLO → rastgele CHALLENGE → AUTH; ardından STATE/EVENT/ACK/HEARTBEAT. IP eşleşmesi anahtardan önce, kimlik ve MAC her çerçevede kontrol edilir. Yeniden bağlantı yeni nonce ve farklı lease alır. Her yönde artan bağlantı sırası, eski/tekrar çerçeveyi reddeder. Çerçeve okuması 2 saniye, yazımı 1 saniye ile sınırlı; aynı anda en fazla 8 bekleyen bağlantı, cihaz başına saniyede 500 mesaj sınırı vardır. Tek cihazın kopması diğerini kapatmaz. Kimlik başına yeni doğrulanmış bağlantı eskisini devralır.

Sunucu gerçek Arena durumundaki değişiklikleri 20 ms aralıklarla kontrol eder; değişiklik yoksa yaklaşık 200 ms'de bir STATE iletir. Firmware heartbeat'i 500 ms'dir ve iletişim 1 saniye kesilirse yakalama bağlamını kapatır. STATE; MatchId, benzersiz oynanış UUID'si, faz/puanlama vardiyası, started/AUTO/TELEOP, fiziksel HUB aktifliği, scoring-enabled, aborted/completed ve Manual Mode bayraklarını taşır. Arduino yalnız son doğrulanmış durumu saklar; ayrı maç süresi hesaplamaz.

Lease yaklaşık 500–600 ms'de veya faz/HUB/puanlama durumu değiştiğinde yenilenir; azami geçerliliği 750 ms'dir. Kabul anında Arena durumu yeniden kontrol edilir. Eski lease, eski faz veya yanlış maça ait yeni olay reddedilir. Gecikmenin fiziksel olay zamanını belirsizleştirdiği yerde önceki faza tahminî puan yazılmaz. Bu tercih bazı sınır/kesinti olaylarının manuel düzeltilmesini gerektirebilir. İlk STATE alınmadan ve uzun bağlantı kesintisinde fiziksel sayım otomatik skora taşınmaz.

Kabul edilen olayın kimliği cihaz+boot+olay sırası+oynanış bağlamından üretilir. Ham sayaç, checkpoint ve tekrarlama günlüğü mevcut manuel işlem transaction'ında atomik yazılır, `Sync` tamamlanmadan başarı ACK'i verilmez. Aynı olay yeniden bağlantı/restart sonrasında ikinci kez sayılmaz. Yazma belirsizliği yeni FUEL işlemlerini kapatır; kurtarma gerekir. Henüz kabul edilmemiş eski maç olayı yeni maça taşınamaz. Kaydedilmiş eski olayın tekrarı yalnız önceki kabulü uzlaştırır.

Sunucu yeniden başladığında mevcut güvenli PostMatch kurtarma devam eder; robotlar/saat/motorlar otomatik başlamaz ve sensör puanlaması açılmaz. Sonuç onayı, aktif tam skor düzenlemesi ve DB bakım işlemleri mevcut koordinatörle sıralanır. Kabul edilmiş sensör olayları manuel işlemlerle aynı günlüğe kaynak/cihaz/boot/olay kimliğiyle eklenir.

## Simülasyon

Yönetim ekranındaki **Isolated sensor simulation** ayrı RAM durumudur; canlı skor, checkpoint, sonuç DB'si, cihaz bağlantısı veya sıralamaya erişen bir yazma yolu yoktur. RED'in AUTO kaybettiği / BLUE'nun kazandığı sabit fixture ile AUTO, transition, Shift1/2, endgame ve post-match grace seçilir. Aynı `Hub.UpdateState` / `Score.Summarize` kuralları kullanılır. Aktif/pasif vardiyadaki sonuç farkı görülebilir. Bu bölüm gerçek maç süresini veya robotları kontrol etmez.

## Testler ve fiziksel kabul

Windows doğrulama komutları:

```powershell
go test ./...
node --test static/js/*test.js
go build -o cheesy-arena-next.exe .
.\tools\Build-Arduino.ps1
```

Otomatik kapsam:

- RED/BLUE ve eşzamanlı iki TCP cihazı; manuel ve sensör transaction'larının eşzamanlılığı.
- Bağlantı kopması, diğer cihazın devamı, reconnect, yeni boot, aynı olayın tekrarı, süreç yeniden açılışı.
- Eski maç ve aynı MatchId'nin yeni oynanışı; eski durum lease'i, faz geçişi, HUB aktif/pasif, maç sonu/onay kilidi.
- Yanlış ID/IP/anahtar, bilinmeyen IP, çerçeve MAC bozulması, bağlantı sırası tekrarı; parçalanmış ve birleşmiş TCP yazımları.
- Manuel düzeltmenin korunması, PLC+Arduino sahipliği, karışık RED/BLUE kaynaklar, Manual Mode.
- Kaynak yapılandırmasının canlı maçta reddi, eski metadata varsayılanı, güvenli kurtarma.
- Simülasyonun canlı skora dokunmaması, admin erişimi ve bilinmeyen simülasyon alanlarının reddi.
- Statik Gi0/23–24 VLAN 100 ve mevcut dinamik switch komutlarının bu alanlara dokunmaması.

İlk tam çalışmada yeni web testindeki yönlendirme beklentisi 303 yerine mevcut sistemin 307 davranışıyla düzeltildi; test kapatılmadı. Son doğrulamada tam Go test paketi ve 17/17 Node testi geçti. RED ve BLUE Uno firmware'leri derlendi: varsayılan tek kanal için her biri **18.348 bayt flash (%56), 986 bayt statik RAM (%48)**; derleyici stack/yerel değişkenler için 1062 bayt kaldığını raporladı. İki kanallı yapılandırma da derlendi (18.392 bayt flash, aynı statik RAM). Bu değerler runtime stack/soak testinin yerine geçmez. Arduino CLI 1.5.2-rc.1, AVR core 1.8.6, Ethernet 2.0.2 ve Crypto 0.4.0 kullanıldı.

`--warnings all` ile Arduino core/Crypto'nun mevcut unused-parameter uyarıları görüldü; derlemeler başarılı. `go vet -composites=false -hostport=false -structtag=false ./...` geçti; tam vet'teki önceden mevcut uyarılar temizlenmiş gibi sunulmuyor. Önceki sürümde belgelenen Windows ThreadSanitizer bellek ayırma sorunu devam eden ortam sınırıdır; race-clean iddiası yoktur.

Çalışan yerel sunucuda ayrıca panel/asset HTTP 200, iki ayrı kapalı cihaz yapılandırması ve simülasyonun canlı skor ile işlem günlüğünü değiştirmediği doğrulandı. Özel anahtarların durum JSON'una veya Git'e girmediği kontrol edildi. Güncellemeden önceki DB yedeği `event-before-arduino-integration.db` olarak yerelde tutuldu.

Gerçek Uno, shield, sensör, saha NIC'i ve Cisco bağlantısı bu oturumda erişilebilir değildi. Firmware derlemesi ve gerçek loopback TCP testleri, fiziksel kabulün yerine geçmez. Sahada her ittifak için ID/IP kontrolü, tek/çift top ve maksimum hız, debounce/pulse aralığı, kuyruk taşması, kablo çekme/yeniden takma, reboot, AUTO/TELEOP ve aktif/pasif grace sınırları, E-stop/A-stop, sonuç onayı ve ekran/TBA çıktısı doğrulanmalıdır.

## Değişen dosyalar

| Alan | Dosyalar |
| --- | --- |
| Ağ protokolü | `fueldevice/config.go`, `protocol.go`, `server.go`, `protocol_test.go` |
| Skor entegrasyonu | `field/fuel_devices.go`, `fuel_device_simulation.go`, `fuel_devices_test.go`; `arena.go`, `manual_fuel.go` |
| Ortak model/kurallar | `model/fuel.go`, `game/hub.go` |
| Başlatma/yapılandırma | `main.go`, `arduino.example.json`, `.gitignore` |
| Panel | `web/arduino.go`, `arduino_test.go`, `web.go`, `templates/arduino.html`, `base.html`, `static/js/arduino.js` |
| Mevcut manuel panel | `templates/manual_fuel.html`, `static/js/manual_fuel.js` kaynak etiketleri |
| Firmware ve araçlar | `firmware/FuelEthernet/FuelEthernet.ino`, `tools/New-ArduinoConfig.ps1`, `Build-Arduino.ps1` |
| Switch regresyonu | `network/switch_test.go`; switch konfigürasyonu değiştirilmedi |

Ekran görüntüleri bu oturumda üretilemedi: otomatik onay denetimi Chrome ile headless görüntü alma komutunu `blocked by policy` gerekçesiyle reddetti. Computer Use alternatifi `native pipe unavailable`, yalıtılmış tarayıcı alternatifi `Browser is not available: iab` döndü. Bu nedenle görsel inceleme yapılmış gibi raporlanmıyor. HTTP sayfa/asset ve yönetim API testleri geçti; arayüz Bootstrap'ın masaüstünde yan yana, tablette alt alta gelen kart düzenini kullanır. [Bağlantı şeması SVG](wiring.svg) kaynak olarak teslim edildi.
