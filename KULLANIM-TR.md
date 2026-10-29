# Cheesy Arena

Bu klasorde Team254/cheesy-arena kaynak kodu ve Windows icin derlenmis uygulama bulunur.

- Baslatmak icin `Baslat.cmd` dosyasina cift tiklayin. Uygulama arka planda baslar ve tarayicida http://localhost:8080/ acilir.
- Kapatmak icin `Durdur.cmd` dosyasina cift tiklayin. Tarayiciyi kapatmak sunucuyu durdurmaz.
- Etkinlik verileri `event.db` dosyasinda saklanir. Yedek almak icin uygulamayi durdurup bu dosyayi kopyalayin.
- Baslatma kayitlari `arena-stdout.log` ve `arena-stderr.log` dosyalarindadir.
- Yeniden derlemek icin bu klasorde `go build -o cheesy-arena.exe .` komutunu calistirin (once uygulamayi durdurun).

Baslatici yerel kullanim icin `-dev` secenegini kullanir. Gercek robotlarla saha kullanimi icin projenin README.md dosyasindaki ag ve donanim ayarlarini uygulayin; bu durumda baslatici yerine `cheesy-arena.exe` komutunu `-dev` olmadan calistirin.

Kaynak: https://github.com/Team254/cheesy-arena
