package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	openai "github.com/sashabaranov/go-openai"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

const (
	OpenAIAPIKeyEnvVar   = "OPENAI_API_KEY"
	HuggingfaceKeyEnvVar = "HUGGINGFACE_API_KEY"
)

func GetEventHandler(client *whatsmeow.Client, gpt *openai.Client) func(interface{}) {
	return func(evt interface{}) {
		// state := client.IsLoggedIn()

		switch v := evt.(type) {
		case *events.Message:
			var messageBody = v.Message.GetConversation()
			client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
				Conversation: proto.String("pong"),
			})
			if messageBody == "ping" {
				client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
					Conversation: proto.String("pong"),
				})
				// response := GenerateGPTResponse(input, gpt)
				// } else if strings.HasPrefix(messageBody, "complete:") {
			}
		}
	}
}

func GenerateGPTResponse(input string, gpt *openai.Client) string {

	resp, err := gpt.CreateCompletion(
		context.Background(),
		openai.CompletionRequest{
			Model:     openai.GPT3TextDavinci002,
			MaxTokens: 4097,
			Prompt:    input,
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}
	if err != nil {
		fmt.Println("Failed to generate GPT response:", err)
		return ""
	}
	return resp.Choices[0].Text
}

//lint:ignore U1000 Ignore unused function warning
func tryconnect(name *string, wg *sync.WaitGroup) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:"+*name+".db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	// Initialize OpenAI GPT
	openaiAPIKey := os.Getenv(OpenAIAPIKeyEnvVar)
	gpt := openai.NewClient(openaiAPIKey)
	if err != nil {
		panic(err)
	}
	client.EmitAppStateEventsOnFullSync = true
	client.AddEventHandler(GetEventHandler(client, gpt))
	var _w fyne.Window
	var _card *widget.Card = nil
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				GUIAPP._app.SendNotification(fyne.NewNotification("scan qrcode with whatsapp", "whatsaap > ... > linked devices"))
				png, err := qrcode.Encode(evt.Code, qrcode.Highest, 512)
				qrcode_img := canvas.NewImageFromReader(bytes.NewReader(png), "qrcode")
				qrcode_img.FillMode = canvas.ImageFillOriginal
				if err != nil {
					panic(fmt.Errorf("hey my man 0: %v", err))
				}
				if _w != nil {
					_card.Content = qrcode_img
				} else {
					_w, _card = initqrwin(qrcode_img, GUIAPP._app)
					_w.Show()
				}
			} else {
				fmt.Println("Login event:", evt.Event)
				if evt.Event == "success" {
					path, err := os.Executable()
					if err != nil {
						fmt.Println("Failed to get the executable path:", err)
						return
					}
					account := Account{
						Name:         *name,
						DatabasePath: path + "/" + *name + ".db",
					}
					SaveAccount(&account)
					accounts, err := LoadAccounts()
					if err != nil {
						fmt.Println("Failed to get the accounts from db:", err)
					}
					store.Accounts = accounts
					_w.Close()
					showtabs(GUIAPP._win)
				}
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			panic(fmt.Errorf("hey my man 1: %v", err))
		} else {
			showtabs(GUIAPP._win)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

//lint:ignore U1000 Ignore unused function warning
func initqrwin(code *canvas.Image, a fyne.App) (fyne.Window, *widget.Card) {
	w := a.NewWindow("Whatsapp QRcode scanner")
	qrcode := widget.NewCard("Scan the QR code", "Whatsapp app > click ... > linked devices > add device", code)
	w.SetContent(container.NewVBox(
		qrcode,
	))
	return w, qrcode
}
