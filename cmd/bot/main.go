package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/massimo-ua/quill/internal/providers/chat/slack"
)

func main() {
	// Create a context that will be canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		log.Println("Received shutdown signal, gracefully stopping...")
		cancel()
	}()

	// Get configuration from environment variables
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")

	if botToken == "" || appToken == "" {
		log.Fatal("SLACK_BOT_TOKEN and SLACK_APP_TOKEN must be set")
	}

	// Create Slack client
	config := slack.NewConfig(botToken, appToken, signingSecret, true)
	factory := slack.NewFactory(config)
	
	chatProvider, err := factory.CreateChatProvider()
	if err != nil {
		log.Fatalf("Failed to create chat provider: %v", err)
	}

	// Start listening for messages
	messageCh, err := chatProvider.ListenForMessages(ctx)
	if err != nil {
		log.Fatalf("Failed to start listening for messages: %v", err)
	}

	// Process messages
	for {
		select {
		case <-ctx.Done():
			log.Println("Context canceled, shutting down...")
			return
		case msg := <-messageCh:
			log.Printf("Received message from %s: %s", msg.Sender(), msg.Content().Text())
			// Here we would normally forward to the domain services for processing
		}
	}
}