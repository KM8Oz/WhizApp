package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	openai "github.com/sashabaranov/go-openai"
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

var GUIAPP = struct {
	_qrwin   func(string, fyne.App) fyne.Window
	_app     fyne.App
	_mainapp func(fyne.App) fyne.Window
}{
	_qrwin:   initqrwin,
	_app:     initapp(),
	_mainapp: initmainwin,
}

type HuggingFaceResponse struct {
	GeneratedText string `json:"generated_text"`
	conversation  struct {
		generated_responses []string `json:"generated_responses"`
		past_user_inputs    []string `json:"past_user_inputs"`
	} `json:"conversation"`
	warnings []string `json:"warnings"`
}

func GetHuggingFaceResponse(prompt string) (string, error) {
	apiKey := os.Getenv(HuggingfaceKeyEnvVar)
	url := "https://api-inference.huggingface.co/models/microsoft/DialoGPT-medium"
	requestBody, err := json.Marshal(map[string]string{
		"inputs": prompt,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var response HuggingFaceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	println("response.generated_text:", response.GeneratedText)
	if len(strings.Fields(response.GeneratedText)) > 0 {
		return response.GeneratedText, nil
	}

	return "", fmt.Errorf("no response from Hugging Face")
}

func GetEventHandler(client *whatsmeow.Client, gpt *openai.Client) func(interface{}) {
	return func(evt interface{}) {
		state, err := client.IsLoggedIn()

		switch v := evt.(type) {
		case *events.Message:
			var messageBody = v.Message.GetConversation()
			fmt.Println("Message event:", messageBody)
			if messageBody == "ping" {
				client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
					Conversation: proto.String("pong"),
				})
				// } else if strings.HasPrefix(messageBody, "complete:") {
			} else {
				// Extract the command arguments
				// args := strings.Fields(messageBody)[1:]
				// Join the arguments to form the input message for GPT
				// input := strings.Join(args, " ")
				// response := GenerateGPTResponse(input, gpt)
				response, err := GetHuggingFaceResponse(messageBody)
				if err != nil {
					fmt.Printf("ChatCompletion error: %v\n", err)
					return
				}
				if len(response) > 0 {
					client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String(response),
					})
				}
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

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:store.db?_foreign_keys=on", dbLog)
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

	client.AddEventHandler(GetEventHandler(client, gpt))
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				GUIAPP._app.SendNotification(fyne.NewNotification("scan qrcode with whatsapp", "whatsaap > ... > linked devices"))
				w := GUIAPP._qrwin(evt.Code, GUIAPP._app)
				w.SetIcon(theme.FyneLogo())
				GUIAPP._qrwin(evt.Code, GUIAPP._app).Show()
				GUIAPP._app.Run()
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
