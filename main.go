package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

// Logger структура для одновременного вывода в файл и консоль
type Logger struct {
	fileLogger *log.Logger
	console    *log.Logger
	mu         sync.Mutex
}

// NewLogger создает новый логгер
func NewLogger(logFilePath string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
		return nil, fmt.Errorf("не удалось создать директорию для логов: %v", err)
	}

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл лога: %v", err)
	}

	return &Logger{
		fileLogger: log.New(logFile, "", log.LstdFlags),
		console:    log.New(os.Stdout, "", log.LstdFlags),
	}, nil
}

// Write реализует io.Writer интерфейс
func (l *Logger) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	n, err = l.fileLogger.Writer().Write(p)
	if err != nil {
		return n, err
	}

	// Также пишем в консоль
	return l.console.Writer().Write(p)
}

// Printf выводит форматированное сообщение
func (l *Logger) Printf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fileLogger.Printf(format, v...)
	l.console.Printf(format, v...)
}

// Println выводит сообщение с новой строки
func (l *Logger) Println(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fileLogger.Println(v...)
	l.console.Println(v...)
}

func main() {
	// Создаем конфигурацию из аргументов командной строки
	config, err := NewTorrentConfig(os.Args)
	if err != nil {
		fmt.Println("Использование: ./torrent-downloader <download-directory> <torrent-file|magnet-link> [log-directory]")
		fmt.Println("  <download-directory> - директория для загрузки файлов")
		fmt.Println("  <torrent-file|magnet-link> - файл торрента или magnet-ссылка")
		fmt.Println("  [log-directory] - опциональная директория для логов (по умолчанию = download-directory)")
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}

	// Создаем директорию для загрузок, если она не существует
	if err := os.MkdirAll(config.DownloadDir, DirPermissions); err != nil {
		log.Fatalf("Не удалось создать директорию для загрузок: %v", err)
	}

	// Получаем абсолютные пути
	absDownloadDir, err := filepath.Abs(config.DownloadDir)
	if err != nil {
		log.Fatalf("Ошибка получения абсолютного пути для директории загрузок: %v", err)
	}
	config.AbsDownloadDir = absDownloadDir

	absLogDir, err := filepath.Abs(config.LogDir)
	if err != nil {
		log.Fatalf("Ошибка получения абсолютного пути для директории логов: %v", err)
	}
	config.AbsLogDir = absLogDir

	// Инициализируем логгер
	logFilePath := filepath.Join(config.AbsLogDir, DefaultLogFileName)
	config.LogFilePath = logFilePath
	logger, err := NewLogger(logFilePath)
	if err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}

	// Перенаправляем стандартный вывод в логгер
	log.SetOutput(io.MultiWriter(logger, os.Stderr))
	log.SetFlags(log.LstdFlags)

	logger.Printf("Файлы будут загружены в: %s\n", config.AbsDownloadDir)
	logger.Printf("Логи будут сохраняться в: %s\n", config.LogFilePath)
	logger.Printf("Директория логов: %s\n", config.AbsLogDir)
	logger.Println("Запуск торрент-клиента...")

	// Конфигурация клиента
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = config.AbsDownloadDir
	cfg.ListenPort = 0

	// Создаем клиент
	client, err := torrent.NewClient(cfg)
	if err != nil {
		logger.Printf("Ошибка создания клиента: %v", err)
		os.Exit(1)
	}
	defer client.Close()

	logger.Println("Клиент торрент успешно запущен")

	// Добавляем торрент
	var t *torrent.Torrent
	if _, err := os.Stat(config.TorrentSource); err == nil {
		metaInfo, err := metainfo.LoadFromFile(config.TorrentSource)
		if err != nil {
			logger.Printf("Ошибка загрузки файла торрента: %v", err)
			os.Exit(1)
		}
		t, err = client.AddTorrent(metaInfo)
		if err != nil {
			logger.Printf("Ошибка добавления торрента: %v", err)
			os.Exit(1)
		}
	} else {
		t, err = client.AddMagnet(config.TorrentSource)
		if err != nil {
			logger.Printf("Ошибка добавления magnet-ссылки: %v", err)
			os.Exit(1)
		}
	}

	// Ждем получения метаданных
	logger.Println("Получаем информацию о торренте...")
	<-t.GotInfo()
	logger.Printf("Загружаем: %s\n", t.Name())

	// Запускаем загрузку всех файлов
	t.DownloadAll()

	// Обработка прерывания программы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool)

	go func() {
		ticker := time.NewTicker(StatsUpdateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				stats := t.Stats()
				progress := float64(t.BytesCompleted()) / float64(t.Info().TotalLength()) * 100
				speed := stats.BytesReadUsefulData.Int64() / 1024 / 5
				peers := stats.ActivePeers

				msg := fmt.Sprintf("Прогресс: %.2f%%, Скорость: %d KB/s, Пиры: %d", progress, speed, peers)
				logger.Println(msg)

				if t.BytesCompleted() == t.Info().TotalLength() {
					logger.Println("Загрузка завершена!")
					done <- true
					return
				}
			case <-done:
				return
			}
		}
	}()

	select {
	case <-done:
		logger.Println("Загрузка торрента успешно завершена")
	case <-sigChan:
		logger.Println("\nПолучен сигнал прерывания, завершаем работу...")

		if t.BytesCompleted() < t.Info().TotalLength() {
			logger.Println("Удаляем неполностью скачанные файлы...")
			if err := removeIncompleteFiles(client, t, logger, config.AbsDownloadDir); err != nil {
				logger.Printf("Ошибка при удалении файлов: %v", err)
			}
		}

		t.Drop()

		os.Exit(1)
	}
}

func removeIncompleteFiles(client *torrent.Client, t *torrent.Torrent, logger *Logger, basePath string) error {
	info := t.Info()
	if info == nil {
		return fmt.Errorf("информация о торренте недоступна")
	}

	bytesCompleted := t.BytesCompleted()

	for _, file := range info.Files {
		if file.Length > bytesCompleted {
			path := filepath.Join(basePath, filepath.Join(file.Path...))
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("ошибка удаления файла %s: %v", path, err)
			}
			logger.Printf("Удален неполностью скачанный файл: %s\n", path)
		}
	}

	return nil
}
