package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

func main() {
	// Проверяем аргументы командной строки
	if len(os.Args) < 2 {
		fmt.Println("Использование: ./torrent-downloader <torrent-file|magnet-link>")
		os.Exit(1)
	}
	torrentSource := os.Args[1]

	// Конфигурация клиента
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = "./downloads" // Папка для загрузок
	cfg.ListenPort = 0          // Случайный порт

	// Создаем клиент
	client, err := torrent.NewClient(cfg)
	if err != nil {
		log.Fatalf("Ошибка создания клиента: %v", err)
	}
	defer client.Close()

	fmt.Println("Клиент торрент запущен")

	// Добавляем торрент
	var t *torrent.Torrent
	if _, err := os.Stat(torrentSource); err == nil {
		// Это файл .torrent
		metaInfo, err := metainfo.LoadFromFile(torrentSource)
		if err != nil {
			log.Fatalf("Ошибка загрузки файла торрента: %v", err)
		}
		t, err = client.AddTorrent(metaInfo)
		if err != nil {
			log.Fatalf("Ошибка добавления торрента: %v", err)
		}
	} else {
		// Это magnet-ссылка
		t, err = client.AddMagnet(torrentSource)
		if err != nil {
			log.Fatalf("Ошибка добавления magnet-ссылки: %v", err)
		}
	}

	// Ждем получения метаданных
	fmt.Println("Получаем информацию о торренте...")
	<-t.GotInfo()
	fmt.Printf("Загружаем: %s\n", t.Name())

	// Запускаем загрузку всех файлов
	t.DownloadAll()

	// Выводим прогресс
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats := t.Stats()
		fmt.Printf("\rПрогресс: %.2f%%, Скорость: %d KB/s, Пиры: %d", 
			t.BytesCompleted()*100/t.Length(), 
			stats.BytesReadUsefulData.Int64()/1024/5, // Средняя скорость за последние 5 сек
			stats.ActivePeers)

		if t.BytesCompleted() == t.Length() {
			fmt.Println("\nЗагрузка завершена!")
			break
		}
	}
}