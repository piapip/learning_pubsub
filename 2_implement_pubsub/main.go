package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/piapip/implement-pubsub/pubsub"
	"go.uber.org/zap"
)

func main() {
	logger := hclog.Default()

	code := 0
	defer gracefullyShutdown(code, logger)

	quit := make(chan os.Signal, 1)

	go func() {
		defer func() {
			close(quit)
		}()

		var (
			err   error
			input string
		)

		var option = 1

		reader := bufio.NewReader(os.Stdin)

		publisher := pubsub.NewPublisher()

		for err == nil && option > 0 {
			option, err = showMenu(reader)

			switch option {
			case 1:
				fmt.Printf("Give a topic: ")
				input, err = reader.ReadString('\n')
				if err != nil {
					log.Fatalln(err)
				}
				input = strings.ReplaceAll(input, "\n", "")

				subscribers := publisher.ListSubscribers(input)
				fmt.Println(subscribers)
			case 2:
				fmt.Printf("Pick a name: ")
				input, err = reader.ReadString('\n')
				if err != nil {
					log.Fatalln(err)
				}
				input = strings.ReplaceAll(input, "\n", "")

				subscriber := pubsub.NewSubscriber(input)
				err = publisher.AddSubscriber(subscriber)
				if err != nil {
					log.Fatalln(err)
				}

				subscriber.Listen()
			case 3:
				topics := publisher.GetTopics()
				fmt.Println(topics)
				fmt.Printf("Pick a topic name: ")
				topicName, err := reader.ReadString('\n')
				if err != nil {
					log.Fatalln(err)
				}
				topicName = strings.ReplaceAll(topicName, "\n", "")

				subscribers := publisher.ListSubscribers(topicName)
				fmt.Println(subscribers)
				fmt.Printf("Subscriber name: ")
				subscriberName, err := reader.ReadString('\n')
				if err != nil {
					log.Fatalln(err)
				}
				subscriberName = strings.ReplaceAll(subscriberName, "\n", "")

				err = publisher.Subscribe(subscriberName, topicName)
				if err != nil {
					log.Println(err)
				}
			case 4:
				topics := publisher.GetTopics()
				fmt.Println(topics)
				fmt.Printf("Pick a topic name: ")
				topicName, err := reader.ReadString('\n')
				if err != nil {
					log.Fatalln(err)
				}
				topicName = strings.ReplaceAll(topicName, "\n", "")

				fmt.Printf("Write a message: ")
				message, err := reader.ReadString('\n')
				if err != nil {
					log.Fatalln(err)
				}
				message = strings.ReplaceAll(message, "\n", "")

				publisher.Publish(topicName, message)
			default:
				if err != nil {
					log.Fatalln(err)
				} else if option <= 0 {
					fmt.Println("Bye!")
				} else {
					fmt.Println("Pick again!")
				}
			}
		}

		code = 1
		if err != nil {
			logger.Error("there's an error", zap.Error(err))
		}
	}()

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	receivedSignal, ok := <-quit

	if ok {
		logger.Info(fmt.Sprintf("SIGNAL %d received, shutting down gracefully, ...", receivedSignal))
	}
}

// 1 publisher multiple subcribers.
func showMenu(reader *bufio.Reader) (int, error) {
	fmt.Println("=============================")
	fmt.Println("1. Show all subscribers.")
	fmt.Println("2. Create a new subscriber.")
	fmt.Println("3. Subscribe.")
	fmt.Println("4. Publish.")
	fmt.Printf("Choose an option: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return -1, err
	}

	option, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return -1, err
	}

	return option, nil
}

func gracefullyShutdown(code int, logger hclog.Logger) {
	logger.Info("Gracefully shutdown...")

	_, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()

	os.Exit(code)
}
