# Torrent Downloader (Go)

[![Go Version](https://avatars.githubusercontent.com/u/203449207?v=4)](https://golang.org/dl/)
[![License: MIT](https://infostart.ru/upload/iblock/174/1746f104e41d54b4787ad1168f55cd07.png)](https://opensource.org/licenses/MIT)

Простое консольное приложение для загрузки торрентов на Go.

## 📦 Установка

```bash
# Клонировать репозиторий
git clone https://github.com/yourusername/torrent-downloader.git
cd torrent-downloader

# Установить зависимости
go get github.com/anacrolix/torrent

# Собрать приложение
go build -o torrent-dl
```

## 🚀 Использование

```bash
# Для .torrent файла
./torrent-dl file.torrent

# Для magnet-ссылки
./torrent-dl "magnet:?xt=urn:btih:..."
```
## 👷 Конфигурация
Измените параметры в коде:
```go
go
cfg.DataDir = "./downloads" // Папка для загрузок
cfg.ListenPort = 0         // Порт (0 = случайный)
```

## 💥 Пример вывода
```text
Downloading: Ubuntu 22.04 LTS
Progress: 54.32%, Speed: 1245 KB/s, Peers: 12
...
Download completed!
```
## 🚑️ Features
+ Supports both torrent files and magnet links

+ Real-time progress reporting

+ Shows download speed and peers

+ Cross-platform (Windows/Linux/macOS)


## ♻️ Лицензия

MIT License. См. [LICENSE](https://opensource.org/licenses/MIT) для подробностей.