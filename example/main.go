package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "modernc.org/sqlite"
	"github.com/mdp/qrterminal/v3"
	"google.golang.org/protobuf/proto"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var client *whatsmeow.Client

func main() {
	// Set up logging
	dbLog := waLog.Stdout("Database", "INFO", true)
	
	// Initialize the database container
	container, err := sqlstore.New(context.Background(), "sqlite", "file:examplestore.db?_pragma=foreign_keys(1)", dbLog)
	if err != nil {
		panic(err)
	}
	
	// Get the first device from the store (or create a new one)
	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		panic(err)
	}
	
	clientLog := waLog.Stdout("Client", "INFO", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	
	// Add event handler
	client.AddEventHandler(eventHandler)
	
	// Connect to WhatsApp
	if client.Store.ID == nil {
		// No ID stored, new login required
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		
		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("QR code:")
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}
	
	// Wait for connection to be established
	fmt.Println("Waiting for connection...")
	time.Sleep(3 * time.Second)
	
	// Example: Uncomment to send a message on startup
	sendExampleMessage()
	
	// Wait for CTRL-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	
	fmt.Println("\nDisconnecting...")
	client.Disconnect()
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Printf("\n📨 New message from %s\n", v.Info.Sender.String())
		
		// Get message text
		messageText := v.Message.GetConversation()
		if messageText == "" {
			messageText = v.Message.GetExtendedTextMessage().GetText()
		}
		
		fmt.Printf("Message: %s\n", messageText)
		
		// Example: Auto-reply to specific messages
		if messageText == "!ping" {
			sendTextMessage(v.Info.Sender, "🏓 Pong!")
		} else if messageText == "!help" {
			helpText := "Available commands:\n!ping - Test bot\n!help - Show this message\n!time - Get current time"
			sendTextMessage(v.Info.Sender, helpText)
		} else if messageText == "!time" {
			currentTime := fmt.Sprintf("⏰ Current time: %s", time.Now().Format("15:04:05 02-Jan-2006"))
			sendTextMessage(v.Info.Sender, currentTime)
		}
		
	case *events.Connected:
		fmt.Println("✅ Connected to WhatsApp!")
		
		// Show your own WhatsApp number (sender)
		ownJID := client.Store.ID
		if ownJID != nil {
			fmt.Printf("🤖 Bot number (sender): %s\n", ownJID.User)
		}
		
	case *events.Disconnected:
		fmt.Println("❌ Disconnected from WhatsApp")
		
	case *events.Receipt:
		if v.Type == events.ReceiptTypeRead {
			fmt.Printf("✓✓ Message %s was read\n", v.MessageIDs[0])
		} else if v.Type == events.ReceiptTypeDelivered {
			fmt.Printf("✓ Message %s was delivered\n", v.MessageIDs[0])
		}
	}
}

// sendTextMessage sends a text message to a recipient
func sendTextMessage(recipient types.JID, text string) {
	msg := &waProto.Message{
		Conversation: proto.String(text),
	}
	
	resp, err := client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Message sent (ID: %s)\n", resp.ID)
}

// sendExampleMessage shows how to send a message to a specific number
func sendExampleMessage() {
	// IMPORTANT: Replace with actual recipient number in international format
	// Format: countrycode + number (no + or spaces)
	// Example: "1234567890" for +1-234-567-890
	recipientNumber := "94779700900" // CHANGE THIS!
	
	// Create JID (WhatsApp ID)
	recipientJID := types.NewJID(recipientNumber, types.DefaultUserServer)
	
	// Create message
	message := &waProto.Message{
		Conversation: proto.String("Hello! This is a test message from whatsmeow bot."),
	}
	
	// Send message
	resp, err := client.SendMessage(context.Background(), recipientJID, message)
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}
	
	fmt.Printf("Message sent successfully! ID: %s, Timestamp: %v\n", resp.ID, resp.Timestamp)
}

// sendMessageWithImage shows how to send an image message
// You'll need to upload the image first using client.Upload()
func sendMessageWithImage(recipient types.JID, imageURL, caption string) {
	msg := &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Caption: proto.String(caption),
			URL:     proto.String(imageURL),
			// Mimetype, FileSize, etc. should be set from the uploaded file
		},
	}
	
	client.SendMessage(context.Background(), recipient, msg)
}

// sendMessageToGroup shows how to send a message to a group
func sendMessageToGroup(groupJID types.JID, text string) {
	msg := &waProto.Message{
		Conversation: proto.String(text),
	}
	
	resp, err := client.SendMessage(context.Background(), groupJID, msg)
	if err != nil {
		fmt.Printf("Error sending group message: %v\n", err)
		return
	}
	
	fmt.Printf("Group message sent! ID: %s\n", resp.ID)
}
