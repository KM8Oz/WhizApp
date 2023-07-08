## WhizApp
Streamline and Empower Your WhatsApp Experience with AI-Driven Management
> WhatsApp Desktop Assistant (GPT3.5-Powered)

Welcome to the WhatsApp Desktop Assistant! This application is a GPT2-powered chatbot that acts as an assistant based on the user's Context . It leverages the power of GPT3.5 turbo, to provide intelligent responses and assist with various tasks.

## Features

- **Intelligent Chatbot**: The assistant is built using GPT3.5turbo, a state-of-the-art language model, which enables it to generate human-like responses based on the user's Context.
- **Google Image optimazed search**: the bot provide an image search engine inside your whatsapp device.
- **Contextual Understanding**: By analyzing the user's Context input from settings, the assistant can understand the context and provide relevant and personalized responses.
- **Task Assistance**: The assistant can perform various tasks such as sending messages, managing targets, add context to the AI assistance .
- **WhatsApp Integration**: The application is seamlessly integrated with WhatsApp, allowing users to interact with their WhatsApp account directly controlled from the desktop application.
- **Golang and Fyne.io Framework**: The application is developed using the Go programming language and utilizes the Fyne.io framework for building the user interface.

## Installation

1. **Clone the repository**:

   ```
   git clone https://github.com/KM8Oz/WhizApp.git
   ```

2. **Install dependencies**:

   ```
   go mod tidy
   ```

3. **Provide WhatsApp API credentials**:

   In order to connect to the WhatsApp service, you need to obtain the necessary API credentials. Obtain the `wpProto` credentials and replace the placeholders in the `config.yaml` file with your actual credentials.

4. **Build and run the application**:

   ```
   go build
   ./whatsapp-desktop-assistant
   ```

## Usage

1. **Launch the application**:

   Once the application is running, the WhatsApp Desktop Assistant will start, and the user interface will be displayed.

2. **Authenticate with WhatsApp**:

   Authenticate the application with your WhatsApp account by scanning the QR code displayed in the application. This will establish a secure connection with your WhatsApp account.

3. **Interact with the Assistant**:

   Begin chatting with the assistant by typing in the input field and pressing Enter. The assistant will generate responses based on your chat history and provide relevant information or assistance.

4. **Perform Tasks**:

   You can ask the assistant to perform various tasks such as sending messages, retrieving chat history, managing contacts, and more. Simply provide the necessary commands or instructions, and the assistant will execute them accordingly.

# Client Documentation

This is the documentation for the `client.go` file in the app. It contains the code related to the Whatsapp client functionality.

## Imports

The following packages are imported in the file:

- `bytes`: Provides functions for working with byte slices.
- `context`: Provides a context for managing concurrent operations.
- `encoding/csv`: Implements CSV encoding and decoding.
- `encoding/json`: Implements JSON encoding and decoding.
- `fmt`: Implements formatted I/O operations.
- `io/ioutil`: Provides I/O utility functions.
- `net/http`: Provides HTTP client and server implementations.
- `os`: Provides a platform-independent interface to operating system functionality.
- `os/signal`: Provides signal handling functionality.
- `os/user`: Provides user account information.
- `path`: Implements utility functions for file paths.
- `strconv`: Implements string conversions.
- `strings`: Implements string manipulation functions.
- `sync`: Provides basic synchronization primitives.
- `syscall`: Provides low-level operating system primitives.
- `time`: Provides functionality for measuring and displaying time.
- `fyne.io/fyne/v2`: Fyne toolkit for building graphical user interfaces.
- `fyne.io/fyne/v2/canvas`: Implements graphical primitives and controls for Fyne.
- `fyne.io/fyne/v2/container`: Implements container types for Fyne.
- `fyne.io/fyne/v2/widget`: Implements UI widgets for Fyne.
- `github.com/mzbaulhaque/gois/pkg/scraper/services`: Package for web scraping services.
- `github.com/sashabaranov/go-openai`: Go client for the OpenAI API.
- `github.com/skip2/go-qrcode`: Library for encoding QR codes.
- `github.com/tmc/langchaingo/chains`: Chain-based implementation for language models.
- `github.com/tmc/langchaingo/llms/openai`: OpenAI-based language model for LLM.
- `github.com/tmc/langchaingo/schema`: Implements schema-related functionality.
- `go.mau.fi/whatsmeow`: Custom package for working with Whatsapp.
- `go.mau.fi/whatsmeow/binary/proto`: Protocol buffer definitions for Whatsapp.
- `go.mau.fi/whatsmeow/store/sqlstore`: SQL-based storage for Whatsapp data.
- `go.mau.fi/whatsmeow/types`: Types and structures related to Whatsapp.
- `go.mau.fi/whatsmeow/types/events`: Event types for Whatsapp.
- `go.mau.fi/whatsmeow/util/log`: Logging utilities for Whatsapp.
- `google.golang.org/protobuf/proto`: Protocol buffer runtime package.

## Functions

The file contains several helper functions:

- `contains`: Checks if a string slice contains a specific item.
- `ConvertToFlickrResult`: Converts data from an interface to a `GoogleResult` struct.
- `askdocument`: Sends a question to the OpenAI language model using previously stored documents.
- `ConvertCSVToDocs`: Converts CSV file content into a slice of `schema.Document` objects.
- `SplitToTokens`: Splits a string into tokens based on a maximum token size.
- `GetTextFormatFromCSV`: Converts CSV file bytes to text format and removes empty rows.
- `GetImageBytes`: Retrieves the byte data and MIME type from an image URL.
- `analyzeCSVData`: Analyzes the CSV data using ChatGPT 3.5 and returns a summary.
- `GetEventHandler`: Returns an event handler function for processing incoming Whatsapp messages.
- `GenerateGPTResponse`: Generates a response using the OpenAI GPT-3.5 language model.

## Additional Functions

There are additional functions in the file for handling QR code scanning, user authentication, and managing the application window.

## Usage

The `client.go` file is an integral part of the Whatsapp Business Manager app. It handles the communication with the Whatsapp server, processes incoming messages, and generates responses using the OpenAI GPT-3.5 language model.

## Dependencies

The file has dependencies on the following packages:

- `fyne.io/fyne/v2`: Fyne toolkit for building graphical user interfaces.
- `github.com/mzbaulhaque/gois/pkg/scraper/services`: Package for web scraping services.
- `github.com/sashabaranov/go-openai`: Go client for the OpenAI API.
- `github.com/skip2/go-qrcode`: Library for encoding QR codes.
- `github.com/tmc/langchaingo/chains`: Chain-based implementation for language models.
- `github.com/tmc/langchaingo/llms/openai`: OpenAI-based language model for LLM.
- `github.com/tmc/langchaingo/schema`: Implements schema-related functionality.
- `go.mau.fi/whatsmeow`: Custom package for working with Whatsapp.
- `go.mau.fi/whatsmeow/binary/proto`: Protocol buffer definitions for Whatsapp.
- `go.mau.fi/whatsmeow/store/sqlstore`: SQL-based storage for Whatsapp data.
- `google.golang.org/protobuf/proto`: Protocol buffer runtime package.


## Configuration

- `wpProto`:
  - `sessionFile`: The file path to store the session information.

## Contributing

Contributions are welcome! If you would like to contribute to this project, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make the necessary changes and commit your code.
4. Push your branch to your forked repository.
5. Submit a pull request with a detailed description of your changes.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgements

- This project utilizes the power of GPT2, a language model developed by OpenAI.
- The application is built using the Go programming language and the Fyne.io framework.
- Special thanks to the open-source community for providing valuable resources and libraries.
