## WhizApp
Streamline and Empower Your WhatsApp Experience with AI-Driven Management
> WhatsApp Desktop Assistant (GPT2-Powered)

Welcome to the WhatsApp Desktop Assistant! This application is a GPT2-powered chatbot that acts as an assistant based on the user's chat history. It leverages the power of GPT2, fine-tuned using machine learning techniques, to provide intelligent responses and assist with various tasks.

## Features

- **Intelligent Chatbot**: The assistant is built using GPT2, a state-of-the-art language model, which enables it to generate human-like responses based on the user's chat history.
- **Contextual Understanding**: By analyzing the user's chat history, the assistant can understand the context and provide relevant and personalized responses.
- **Task Assistance**: The assistant can perform various tasks such as sending messages, retrieving chat history, managing contacts, and more.
- **WhatsApp Integration**: The application is seamlessly integrated with WhatsApp, allowing users to interact with their WhatsApp account directly from the desktop application.
- **Golang and Fyne.io Framework**: The application is developed using the Go programming language and utilizes the Fyne.io framework for building the user interface.

## Installation

1. **Clone the repository**:

   ```
   git clone https://github.com/your/repository.git
   ```

2. **Install dependencies**:

   ```
   go mod download
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

## Configuration

The application can be configured using the `config.yaml` file. The following configuration options are available:

- `wpProto`:
  - `phoneNumber`: Your WhatsApp account phone number.
  - `clientName`: The name of the client.
  - `sessionFile`: The file path to store the session information.
  - `timeout`: The timeout duration for the WhatsApp connection.

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
