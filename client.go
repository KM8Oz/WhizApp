package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mzbaulhaque/gois/pkg/scraper/services"
	openai "github.com/sashabaranov/go-openai"
	"github.com/skip2/go-qrcode"
	"github.com/tmc/langchaingo/chains"
	openailc "github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var _req = map[string]openai.ChatCompletionRequest{}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
func ConvertToFlickrResult(data interface{}) (services.GoogleResult, bool) {
	result := services.GoogleResult{}

	// Use reflection to access the fields of the interface
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling interface:", err)
		return result, false
	}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return result, false
	}
	return result, true
}
func askdocument(gpt *openai.Client, command string, user string) (string, error) {
	llm, err := openailc.New(openailc.WithToken(openaiAPIKey))
	if err != nil {
		return "", err
	}
	if _docs, ok := globaldocs[user]; ok {
		globaldocs[user] = _docs
		// Prompt the LLM using the docs
		stuffQAChain := chains.LoadStuffQA(llm)
		result, err := chains.Call(context.Background(), stuffQAChain, map[string]interface{}{
			"input_documents": _docs,
			"question":        command,
		})
		if err != nil {
			return "", err
		}
		if len(result) > 0 {
			output := result["text"].(string)
			return output, nil
		}
		return "", nil
	} else {
		return "", fmt.Errorf("you have no document in memory to QA!")
	}
}

// ConvertCSVToDocs converts CSV file content into a slice of schema.Document objects,
// where each document represents a page of 100 rows from the CSV file.
func ConvertCSVToDocs(csvContent string) ([]schema.Document, error) {
	records := SplitToTokens(csvContent)
	var docs []schema.Document = []schema.Document{}
	for _, pageContent := range records {
		docs = append(docs, schema.Document{PageContent: pageContent})
	}

	return docs, nil
}
func SplitToTokens(prompt string) []string {
	const maxTokens = 2000
	var tokens []string

	// Split the prompt into individual words
	words := strings.Fields(prompt)

	// Iterate over the words and combine them into tokens
	var tokenBuilder strings.Builder
	for _, word := range words {
		if tokenBuilder.Len()+len(word) > maxTokens {
			// Token size exceeds the limit, append the current token and start a new one
			tokens = append(tokens, tokenBuilder.String())
			tokenBuilder.Reset()
		}
		tokenBuilder.WriteString(word)
		tokenBuilder.WriteByte(' ')
	}

	// Append the last token
	tokens = append(tokens, tokenBuilder.String())

	return tokens
}

// GetTextFormatFromCSV converts CSV file bytes to text format and removes empty rows
func GetTextFormatFromCSV(csvData []byte) (string, error) {
	reader := csv.NewReader(bytes.NewReader(csvData))
	lines, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	var output strings.Builder
	for _, line := range lines {
		// Skip empty rows
		if len(line) == 0 {
			continue
		}

		// Remove empty cells in a row
		var nonEmptyCells []string
		for _, cell := range line {
			if cell != "" {
				nonEmptyCells = append(nonEmptyCells, cell)
			}
		}

		// Append non-empty cells in a row
		if len(nonEmptyCells) > 0 {
			output.WriteString(strings.Join(nonEmptyCells, ","))
			output.WriteString("\n")
		}
	}

	return output.String(), nil
}
func GetImageBytes(url string) ([]byte, string, error) {
	// Send a GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	// Check if the response status code is OK
	if response.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP request failed with status code %d", response.StatusCode)
	}

	// Read the response body
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	// Check if the URL points to an image
	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, "", fmt.Errorf("URL does not point to an image")
	}

	return bodyBytes, contentType, nil
}

// analyzeCSVData analyzes the CSV data using ChatGPT 3.5 and returns a summary
func analyzeCSVData(csvData string, gpt *openai.Client, command string, user string) (string, error) {
	llm, err := openailc.New(openailc.WithToken(openaiAPIKey))
	if err != nil {
		return "", err
	}
	_docs, err := ConvertCSVToDocs(csvData)
	if err != nil {
		return "", err
	}
	globaldocs[user] = _docs[:1]
	// Prompt the LLM using the docs
	stuffQAChain := chains.LoadStuffQA(llm)
	result, err := chains.Call(context.Background(), stuffQAChain, map[string]interface{}{
		"input_documents": globaldocs[user],
		"question":        command,
	})
	if err != nil {
		return "", err
	}
	if len(result) > 0 {
		output := result["text"].(string)
		return output, nil
	}
	return "", nil
}
func GetEventHandler(client *whatsmeow.Client, gpt *openai.Client) func(interface{}) {
	questions := []string{
		"How can I assist you today?",
		"What specific information are you looking for?",
		"Is there a particular feature you need help with?",
		"Do you have any technical issues that need troubleshooting?",
	}
	// Create a button template with the predefined questions
	QuickReplybuttons := make([]*waProto.HydratedTemplateButton, len(questions))
	QuickReplybuttons_ := make([]*waProto.ButtonsMessage_Button, len(questions))
	for i, question := range questions {
		QuickReplybuttons_[i] = &waProto.ButtonsMessage_Button{
			ButtonId: proto.String(strconv.Itoa(i)),
			ButtonText: &waProto.ButtonsMessage_Button_ButtonText{
				DisplayText: proto.String(question),
			},
		}
	}
	for i, question := range questions {
		QuickReplybuttons[i] = &waProto.HydratedTemplateButton{
			Index: proto.Uint32(uint32(i)),
			HydratedButton: &waProto.HydratedTemplateButton_QuickReplyButton{
				QuickReplyButton: &waProto.HydratedTemplateButton_HydratedQuickReplyButton{
					DisplayText: proto.String(question),
					Id:          proto.String(question),
				},
			},
		}
	}
	hydratedCallButton := &waProto.HydratedTemplateButton{
		Index: proto.Uint32(uint32(10)),
		HydratedButton: &waProto.HydratedTemplateButton_CallButton{
			CallButton: &waProto.HydratedTemplateButton_HydratedCallButton{
				DisplayText: proto.String("Call US"),
				PhoneNumber: proto.String("+212709251456"),
			},
		},
	}
	QuickReplybuttons = append(QuickReplybuttons, hydratedCallButton)
	// buttons_title := "Please select one of the following questions:"
	// hydratedFourRowTemplate := waProto.TemplateMessage_HydratedFourRowTemplate{
	// 	HydratedContentText: proto.String("الآن العرض الجديد"),
	// 	HydratedFooterText:  proto.String("تطبّق الشروط والأحكام"),
	// 	HydratedButtons:     QuickReplybuttons,
	// 	TemplateId:          proto.String("id1"),
	// 	Title: &waProto.TemplateMessage_HydratedFourRowTemplate_HydratedTitleText{
	// 		HydratedTitleText: buttons_title,
	// 	},
	// }
	// templateMessage := waProto.TemplateMessage{
	// 	// ContextInfo:      &waProto.ContextInfo{},
	// 	HydratedTemplate: &hydratedFourRowTemplate,
	// 	TemplateId:       proto.String("kom"),
	// 	Format:           &waProto.TemplateMessage_FourRowTemplate_{},
	// }
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.LoggedOut:
			err := client.Connect()
			if err != nil {
				panic(err)
			}
		case *events.Message:
			var messageBody = v.Message.GetConversation()
			fmt.Println("Message event:", v.Message.GetConversation(), v.Info.Type)
			client.MarkRead([]string{v.Info.ID}, time.Now(), v.Info.Chat, v.Info.Sender)
			switch {
			case v.IsDocumentWithCaption:
				DocumentWithCaption := v.Message.DocumentMessage
				if bytes, _error := client.Download(DocumentWithCaption); _error == nil {
					switch DocumentWithCaption.GetMimetype() {
					case "text/csv":
						if csvfile, csvfileerr := GetTextFormatFromCSV(bytes); csvfileerr == nil {
							if res, err := analyzeCSVData(csvfile, gpt, DocumentWithCaption.GetCaption(), v.Info.Sender.String()); err == nil {
								client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
									Conversation: proto.String(res),
								})
							} else {
								client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
									Conversation: proto.String(fmt.Sprintf("__%s__", err.Error())),
								})
							}
						} else {
							client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
								Conversation: proto.String(fmt.Sprintf("__%s__", csvfileerr.Error())),
							})
						}
					default:
						client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
							Conversation: proto.String("File format is not implimented yet!"),
						})
					}
				}
			case v.Info.Type == "media" && !v.IsDocumentWithCaption:
				client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
					Conversation: proto.String("File format is not implimented yet!"),
				})
			case strings.ToLower(messageBody) == "ping":
				client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
					Conversation: proto.String("pong"),
				})
			case strings.HasPrefix(strings.ToLower(messageBody), "/askdoc"):
				args := strings.Fields(messageBody)[1:]
				the_rest := strings.Join(args, " ")
				if res, err := askdocument(gpt, the_rest, v.Info.Sender.String()); err == nil {
					client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String(res),
					})
				} else {
					client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String(fmt.Sprintf("__%s__", err.Error())),
					})
				}
			case strings.HasPrefix(strings.ToLower(messageBody), "/reset"):
				if _, ok := _req[v.Info.Sender.String()]; ok {
					var _allmessages []openai.ChatCompletionMessage = []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleSystem,
							Content: global_context,
						},
					}
					_req[v.Info.Sender.String()] = openai.ChatCompletionRequest{
						Model:    openai.GPT3Dot5Turbo,
						Messages: _allmessages,
					}
				}
			case strings.HasPrefix(strings.ToLower(messageBody), "/new"):
				args := strings.Fields(messageBody)[1:]
				the_rest := strings.Join(args, " ")
				if the_rest == "" {
					var _allmessages []openai.ChatCompletionMessage = []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleSystem,
							Content: global_context,
						},
					}
					_req[v.Info.Sender.String()] = openai.ChatCompletionRequest{
						Model:    openai.GPT3Dot5Turbo,
						Messages: _allmessages,
					}
				} else {
					var _allmessages []openai.ChatCompletionMessage = []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleSystem,
							Content: the_rest,
						},
					}
					_req[v.Info.Sender.String()] = openai.ChatCompletionRequest{
						Model:    openai.GPT3Dot5Turbo,
						Messages: _allmessages,
					}
				}
			// case strings.HasPrefix(strings.ToLower(messageBody), "/set_group_name"):
			// 	args := strings.Fields(messageBody)[1:]
			// 	name := strings.Join(args, " ")
			// 	if irr := client.SetGroupName(v.Info.Chat, name); irr != nil {
			// 		_, err := client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
			// 			Conversation: proto.String(irr.Error()),
			// 		})
			// 		if err != nil {
			// 			fmt.Printf("ImageMessage error: %v\n", err)
			// 		}
			// 	}
			case strings.HasPrefix(strings.ToLower(messageBody), "/image"):
				args := strings.Fields(messageBody)[1:]
				query := strings.Join(args, " ")
				fmt.Printf("query: %s", query)
				// if up, err := client.Upload(context.Background(), bytedata, whatsmeow.MediaImage); err != nil {
				// 	return nil, err
				//   } else {

				// 	var message = &waProto.ImageMessage{
				// 	  Url:           &up.URL,
				// 	  Mimetype:      proto.String(mimetype),
				// 	  Caption:       proto.String("Caption"),
				// 	  FileSha256:    up.FileSHA256,
				// 	  FileEncSha256: up.FileEncSHA256,
				// 	  FileLength:    &up.FileLength,
				// 	  MediaKey:      up.MediaKey,
				// 	  DirectPath:    &up.DirectPath,
				// 	}
				//   }
				config := &services.GoogleConfig{
					Query: query,
				}
				gs := &services.GoogleScraper{Config: config}
				items, _, _err := gs.Scrape()
				if _err != nil {
					fmt.Printf("ImageMessage error: %v\n", _err)
					client.SendPresence(types.PresenceAvailable)
					_, err := client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
						Conversation: proto.String("images not found!"),
					})
					if err != nil {
						fmt.Printf("ImageMessage error: %v\n", err)
					}
				}
				var convertedItems = make([]services.GoogleResult, 4)
				for i, item := range items {
					if i <= 3 {
						data, ok := ConvertToFlickrResult(item)
						if ok {
							convertedItems[i] = data
						}
					} else {
						break
					}
				}
				fmt.Print(convertedItems)
				if len(convertedItems) > 0 {
					for i := 0; i < len(convertedItems); i++ {
						if convertedItems[i].URL != "" {
							bytedata, mimeType, __err := GetImageBytes(convertedItems[i].URL)
							if __err != nil {
								fmt.Printf("ImageMessage error: %v\n", __err)
								return
							} else {
								if up, err := client.Upload(context.Background(), bytedata, whatsmeow.MediaImage); err == nil {
									var message = &waProto.ImageMessage{
										Url:           &up.URL,
										Mimetype:      proto.String(mimeType),
										Caption:       proto.String(convertedItems[i].Title),
										FileSha256:    up.FileSHA256,
										FileEncSha256: up.FileEncSHA256,
										FileLength:    &up.FileLength,
										MediaKey:      up.MediaKey,
										DirectPath:    &up.DirectPath,
									}
									_, err := client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
										ImageMessage: message,
									})
									if err != nil {
										fmt.Printf("ImageMessage error: %v\n", err)
										return
									}
								} else {
									fmt.Printf("Upload error: %v", err)
								}
							}
						}
					}

				}
			default:
				if !v.Info.Sender.IsEmpty() {
					fmt.Printf("new message: %s in %s\n", v.Info.Sender.String(), v.Info.Chat.String())
					_, err := infoTargetsbng.GetValue(v.Info.Sender.User)
					istargeted := len(infoTargetsbng.Keys()) == 0 || err == nil
					GroupInfo, err := client.GetGroupInfo(v.Info.Chat)
					isgrouptargeted := targetedgroups == "all" || GroupInfo.Name == targetedgroups
					if istargeted && isgrouptargeted {
						response, err := GenerateGPTResponse(messageBody, v.Info.Sender.String(), gpt)
						// // response, err := GetHuggingFaceResponse(messageBody)
						if err != nil {
							fmt.Printf("ChatCompletion error: %v\n", err)
							return
						}
						if len(response) > 0 {
							// Create a buttons message.
							client.SendPresence(types.PresenceAvailable)
							_, err := client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
								Conversation: proto.String(response),
							})
							if err != nil {
								fmt.Printf("ERROR Message: %v", err)
							}
						}
					}
				}
			}
		}
	}
}

func GenerateGPTResponse(input string, user string, gpt *openai.Client) (string, error) {
	var _allmessages []openai.ChatCompletionMessage
	if _, ok := _req[user]; !ok {

		_allmessages = []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: global_context,
			},
		}
	}
	_allmessages = append(_allmessages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})
	_req[user] = openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: _allmessages,
	}
	resp, err := gpt.CreateChatCompletion(
		context.Background(),
		_req[user],
	)
	if err != nil {
		return "fails!!!", fmt.Errorf("chatCompletion error: %v", err)
	}
	_allmessages = append(_allmessages, resp.Choices[0].Message)
	_req[user] = openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: _allmessages,
	}
	return resp.Choices[0].Message.Content, nil
}
func getUserHomeDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return currentUser.HomeDir, nil
}

//lint:ignore U1000 Ignore unused function warning
func tryconnect(name *string, wg *sync.WaitGroup) {
	__path, err := getUserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	dbpath := path.Join(__path, *name+".db")
	container, err := sqlstore.New("sqlite3", "file:"+dbpath+"?_foreign_keys=on", dbLog)
	if err != nil {
		fmt.Println(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		fmt.Println(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	// Initialize OpenAI GPT
	if err != nil {
		fmt.Println(err)
	}
	GUIAPP.client.EnableAutoReconnect = true
	client.EmitAppStateEventsOnFullSync = true
	var _w fyne.Window
	var _card *widget.Card = nil
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			fmt.Println(err)
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
					if err != nil {
						fmt.Println("Failed to get the executable path:", err)
						return
					}
					account := Account{
						Name:         *name,
						DatabasePath: dbpath,
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
