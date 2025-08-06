# Torrent Downloader (Go)

[![Go Version](https://avatars.githubusercontent.com/u/203449207?v=4)](https://golang.org/dl/)
[![License: MIT](https://infostart.ru/upload/iblock/174/1746f104e41d54b4787ad1168f55cd07.png)](https://opensource.org/licenses/MIT)

Простое консольное приложение для загрузки торрентов на Go с поддержкой логирования и гибкой конфигурации.

## 📦 Установка

```bash
# Клонировать репозиторий
git clone https://github.com/yourusername/torrent-downloader.git
cd torrent-downloader

# Установить зависимости
go get github.com/anacrolix/torrent

# Собрать приложение
go build -o torrent-downloader
```

## 🚀 Использование

```bash
# Для .torrent файла
./torrent-downloader /path/to/download/directory /path/to/torrent/file.torrent

# Для magnet-ссылки
./torrent-downloader /path/to/download/directory "magnet:?xt=urn:btih:..."

# С указанием директории для логов
./torrent-downloader /path/to/download/directory /path/to/torrent/file.torrent /path/to/logs
```

### Параметры командной строки:

- `<download-directory>` - директория для загрузки файлов (обязательный)
- `<torrent-file|magnet-link>` - файл торрента или magnet-ссылка (обязательный)
- `[log-directory]` - опциональная директория для логов (по умолчанию = download-directory)

## 📁 Структура проекта

```
torrent/
├── main.go          # Основная логика приложения
├── config.go        # Конфигурационные переменные и структуры
├── go.mod           # Зависимости Go
├── go.sum           # Хеши зависимостей
├── README.md        # Документация
└── LICENSE.txt      # Лицензия
```

## 🔧 Конфигурация

### Основные настройки (config.go):

```go
// Интервал обновления статистики
StatsUpdateInterval = 5 * time.Second

// Права доступа для директорий
DirPermissions = 0755

// Имя файла лога по умолчанию
DefaultLogFileName = "torrent-downloader.log"
```

### Настройка клиента (main.go):

```go
cfg := torrent.NewDefaultClientConfig()
cfg.DataDir = config.AbsDownloadDir  // Директория для загрузок
cfg.ListenPort = 0                   // Случайный порт
```

## 📊 Логирование

Приложение поддерживает двойное логирование:
- **Консольный вывод** - для мониторинга в реальном времени
- **Файловое логирование** - для сохранения истории

### Логи содержат:
- Информацию о директориях загрузки и логов
- Прогресс загрузки в реальном времени
- Скорость загрузки и количество пиров
- Ошибки и предупреждения
- Статистику завершения

### Пример лога:
```
2025/08/06 15:00:52 Файлы будут загружены в: /path/to/downloads
2025/08/06 15:00:52 Логи будут сохраняться в: /path/to/logs/torrent-downloader.log
2025/08/06 15:00:52 Директория логов: /path/to/logs
2025/08/06 15:00:52 Прогресс: 54.32%, Скорость: 1245 KB/s, Пиры: 12
2025/08/06 15:00:52 Загрузка завершена!
```

## 💥 Пример вывода

```text
Использование: ./torrent-downloader <download-directory> <torrent-file|magnet-link> [log-directory]
  <download-directory> - директория для загрузки файлов
  <torrent-file|magnet-link> - файл торрента или magnet-ссылка
  [log-directory] - опциональная директория для логов (по умолчанию = download-directory)

Файлы будут загружены в: /Users/user/Downloads
Логи будут сохраняться в: /Users/user/Downloads/torrent-downloader.log
Директория логов: /Users/user/Downloads
Запуск торрент-клиента...
Клиент торрент успешно запущен
Получаем информацию о торренте...
Загружаем: Ubuntu 22.04 LTS
Прогресс: 54.32%, Скорость: 1245 KB/s, Пиры: 12
...
Загрузка завершена!
Загрузка торрента успешно завершена
```

## 🚑️ Features

### ✅ Основные возможности:
- Поддержка .torrent файлов и magnet-ссылок
- Отображение прогресса в реальном времени
- Показ скорости загрузки и количества пиров
- Кроссплатформенность (Windows/Linux/macOS)

### ✅ Логирование:
- Двойное логирование (консоль + файл)
- Настраиваемая директория логов
- Подробная информация о процессе загрузки

### ✅ Конфигурация:
- Вынесенные в отдельный файл константы
- Гибкая настройка параметров
- Валидация аргументов командной строки

### ✅ Обработка ошибок:
- Удаление неполностью скачанных файлов при прерывании
- Graceful shutdown при получении сигналов
- Подробные сообщения об ошибках

## 🔧 Разработка

### Структура кода:
- **main.go** - основная логика приложения
- **config.go** - конфигурационные переменные и структуры
- **Logger** - кастомный логгер с двойным выводом

### Добавление новых функций:
1. Добавьте константы в `config.go`
2. Обновите структуру `TorrentConfig` при необходимости
3. Используйте `logger` для вывода информации

## ♻️ Лицензия

MIT License. См. [LICENSE.txt](LICENSE.txt) для подробностей.