package main

import (
	"github.com/taiyoh/go-cm160"
	"log"
	"time"
)

func main() {

	config := LoadConfig()

	recCh := make(chan *cm160.Record)
	sigCh := BuildSigWatcher()

	client := NewMkrClient(config.Mackerel, config.Name)

	device := cm160.Open(recCh)
	defer device.Close()

	log.Println("device initialized")

	go func() {
		for {
			select {
			case <-sigCh:
				log.Println("stop running")
				device.Stop()
			case record := <-recCh:
				if record.IsLive {
					log.Printf("live record amps: %#v\n", record.Amps)
				} else {
					log.Printf("not live at %d-%02d-%02d %02d:%02d amps: %#v\n", record.Year, record.Month, record.Day, record.Hour, record.Minute, record.Amps)
				}
				client.post(record)
				time.Sleep(50 * time.Millisecond)
			default:
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	device.Run()

	log.Println("exit process")
}
