package main

import (
	"bytes"
	"encoding/json"
	"image"
	"log"
	"net/http"
	"os"
	"test/config"
	"test/internal/domain"
	rabbitmq "test/internal/rabbitMq"
	"test/pkg"
	"time"

	"github.com/disintegration/imaging"
)

func connectRabbit(rabbitURL string) (*rabbitmq.RabbitMQConnectionManager, error) {
	var objectReciever *rabbitmq.RabbitMQConnectionManager
	var err error
	for i := 0; i < 15; i++ {
		objectReciever, err = rabbitmq.NewRabbitMQConnection(rabbitURL, "que")
		if err == nil {
			return objectReciever, nil
		}
		log.Printf("RabbitMQ connection failed (%d/15): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Printf("Warning: error loading .env: %v (using defaults)", err)
		cfg, _ = config.Load("")
	}

	objectReciever, err := connectRabbit(cfg.RabbitMQURL())
	if err != nil {
		log.Fatalf("failed creating RabbitMq connection %v", err)
	}
	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	msgs, err := objectReciever.Receive()
	if err != nil {
		log.Fatalf("failed creating RabbitMq coonection %v", err)
	}

	go func() {
		for msg := range msgs {
			var task domain.Task
			err := json.Unmarshal(msg.Body, &task)
			if err != nil {
				continue
			}

			img, err := pkg.FromFileToImage(task.ImageBase)
			if err != nil {
				log.Printf("Failed to decode image %q: %v", task.ImageBase, err)
				// Попробуем относительный путь из рабочего каталога (Docker: /app)
				altPath := "./" + task.ImageBase
				img, err = pkg.FromFileToImage(altPath)
				if err != nil {
					altPath2 := "/app/" + task.ImageBase
					img, err = pkg.FromFileToImage(altPath2)
					if err != nil {
						log.Printf("Also failed on alternate paths: %q, %q: %v", altPath, altPath2, err)
						continue
					}
				}
			}

			var editedImg *image.NRGBA
			switch task.FilterName {
			case "Sharpen":
				filterParam, ok := task.FilterParametes.(float64)
				if !ok {
					log.Printf("Invalid type for FilterParametes, expected float64")
					continue
				}
				editedImg = imaging.Sharpen(img, filterParam)
			case "Invert":
				{
					editedImg = imaging.Invert(img)
				}
			}
			if editedImg == nil {
				log.Printf("No valid filter applied for task ID: %s", task.Id)
				continue
			}

			task.Result = "images/output_image.png"

			outputFile, err := os.Create(task.Result)
			if err != nil {
				log.Printf("Failed to create output file: %v", err)
				continue
			}

			err = imaging.Save(editedImg, task.Result)
			if err != nil {
				log.Printf("Failed to save image: %v", err)
				outputFile.Close()
				continue
			}

			log.Println("Image successfully saved to file")
			outputFile.Close()

			task.Status = "done"

			body, _ := json.Marshal(task)
			commitURL := cfg.AppURL() + "/commit"
			_, err = http.Post(commitURL, "application/json", bytes.NewBuffer(body))
			if err != nil {
				log.Printf("Failed to send result to /commit (%s): %v", commitURL, err)
			}
		}
	}()
	select {}
}
