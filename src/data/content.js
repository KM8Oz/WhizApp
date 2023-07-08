const content = {
  main: {
    title: "WhizApp",
    description:
      "Streamline and Empower Your WhatsApp Experience with AI-Driven Management",
    intro:
      "Welcome to the WhatsApp Desktop Assistant! This application is a GPT2-powered chatbot that acts as an assistant based on the user’s chat history. It leverages the power of GPT2, fine-tuned using machine learning techniques, to provide intelligent responses and assist with various tasks.",
    logo: "icon.png",
    repository: "https://github.com/KM8Oz/WhizApp",
    link: "http://dl.whizapp.dup.company",
  },
  Features: [
    {
      feature: "Intelligent Chatbot",
      intro:
        "The assistant is built using GPT2, a state-of-the-art language model, which enables it to generate human-like responses based on the user’s chat history.",
    },

    {
      feature: "Contextual Understanding",
      intro:
        "By analyzing the user's chat history, the assistant can understand the context and provide relevant and personalized responses.",
    },
    {
      feature: "Task Assistance",
      intro:
        "The assistant can perform various tasks such as sending messages, retrieving chat history, managing contacts, and more.",
    },
    {
      feature: "WhatsApp Integration",
      intro:
        "The application is seamlessly integrated with WhatsApp, allowing users to interact with their WhatsApp account directly from the desktop application.",
    },
    {
      feature: "Golang and Fyne.io Framework",
      intro:
        "The application is developed using the Go programming language and utilizes the Fyne.io framework for building the user interface.",
    },
  ],

  Usage: [
    {
      step: "Launch the application",
      description:
        "Once the application is running, the WhatsApp Desktop Assistant will start, and the user interface will be displayed.",
    },
    {
      step: "Authenticate with WhatsApp",
      description:
        "Authenticate the application with your WhatsApp account by scanning the QR code displayed in the application. This will establish a secure connection with your WhatsApp account",
    },
    {
      step: "Interact with the Assistant",
      description:
        "Begin chatting with the assistant by typing in the input field and pressing Enter. The assistant will generate responses based on your chat history and provide relevant information or assistance.",
    },
    {
      step: "Perform Tasks",
      description:
        "You can ask the assistant to perform various tasks such as sending messages, retrieving chat history, managing contacts, and more. Simply provide the necessary commands or instructions, and the assistant will execute them accordingly.",
    },
  ],

  Configuration: [
    {
      config: "phoneNumber",
      intro: "Your WhatsApp account phone number.",
    },{
      config: "clientName",
      intro: "The name of the client.",
    },{
      config: "sessionFile",
      intro: "The file path to store the session information.",
    },{
      config: "timeout",
      intro: "The timeout duration for the WhatsApp connection.", 
    },
  ],

  Acknowledgements: [
    {
      skill: "GPT2",
      description:
        "This project utilizes the power of GPT2, a language model developed by OpenAI",
    },{
      skill: "GoLang",
      description:
        "The application is built using the Go programming language and the Fyne.io framework.",
    },{
      skill: "Github",
      description:
        "Special thanks to the open-source community for providing valuable resources and libraries.",
    },
  ],



  
  
};

export default content;
