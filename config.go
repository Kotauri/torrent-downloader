package main

import (
	"errors"
	"time"
)

// Ошибки конфигурации
var (
	ErrInsufficientArgs = errors.New("недостаточно аргументов")
	ErrTooManyArgs      = errors.New("слишком много аргументов")
)

// Конфигурационные константы
const (
	// Интервал обновления статистики загрузки
	StatsUpdateInterval = 5 * time.Second

	// Права доступа для директорий
	DirPermissions = 0755

	// Права доступа для файлов логов
	LogFilePermissions = 0644

	// Имя файла лога по умолчанию
	DefaultLogFileName = "torrent-downloader.log"

	// Минимальное количество аргументов командной строки
	MinArgs = 3

	// Максимальное количество аргументов командной строки
	MaxArgs = 4
)

// Конфигурация клиента торрента
type TorrentConfig struct {
	// Директория для загрузки файлов
	DownloadDir string

	// Директория для логов
	LogDir string

	// Источник торрента (файл или magnet-ссылка)
	TorrentSource string

	// Абсолютный путь к директории загрузок
	AbsDownloadDir string

	// Абсолютный путь к директории логов
	AbsLogDir string

	// Путь к файлу лога
	LogFilePath string
}

// NewTorrentConfig создает новую конфигурацию из аргументов командной строки
func NewTorrentConfig(args []string) (*TorrentConfig, error) {
	if len(args) < MinArgs {
		return nil, ErrInsufficientArgs
	}

	if len(args) > MaxArgs {
		return nil, ErrTooManyArgs
	}

	config := &TorrentConfig{
		DownloadDir:   args[1],
		TorrentSource: args[2],
	}

	// Опциональный аргумент - директория для логов
	if len(args) >= 4 {
		config.LogDir = args[3]
	} else {
		// По умолчанию логи сохраняются в директорию загрузок
		config.LogDir = config.DownloadDir
	}

	return config, nil
}
