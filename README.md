
# CrossPlatform-TechSupport-Chatbot
A multi-platform chatbot that provides intelligent customer and tech support using OpenAI, Dialogflow, and a Retrieval-Augmented Generation (RAG) system.

## Table of Contents
- [Features](#features)
- [Demo](#demo)
- [How It Works](#how-it-works)
- [Installation](#installation)
- [Tech Stack](#tech-stack)
- [License](#license)

## Features
- Multi-platform support (Messenger, Telegram, LINE, or custom web page)
- Document storing and chunking
- OpenAI integration for conversational AI and text embedding/semantic search
- Retrieval-Augmented Generation (RAG) for document-based responses
- Context-aware responses and smart routing
[ - - Handles FAQs, troubleshooting, and customer inquiries]: #

## Demo
- **Video Demo**: [Coming Soon]
- **Live Demo**: [Coming Soon]

## How It Works
- Users interact with the chatbot through various platforms (support FB Messenger, Telegram, LINE, or custom platform).
- Dialogflow detects intents for common questions and support requests.
- OpenAI generates conversational responses for unrecognized inputs.
- The Retrieval-Augmented Generation (RAG) system fetches relevant documents for FAQs and troubleshooting.
- document upload.

## Tech Stack
- **Frontend**: React (web interaction)
- **Backend**: Go (Gin framework)
- **Database**: PostgreSQL
- **APIs**: OpenAI API, Dialogflow, META APIs, Telegram API, LINE API
[ - **Cloud**: AWS (for deployment) ]: #

## Installation
- **Clone the Repository**:
   ```bash
   git clone https://github.com/petersun1937/CrossPlatform-TechSupport-Chatbot.git
   cd CrossPlatform-TechSupport-Chatbot
   ```

- **Backend Setup**:
   - Ensure Go and Python are installed.
   - Install the required Python package `pdfplumber`:
     ```bash
     pip install pdfplumber
     ```
   - Set up environment variables in the `configs/.env` file (e.g., API keys, database configuration).
   - Run the script by
     ```bash
     go run main.go
     ```


- **Frontend Setup**:
   - Navigate to `React_custom_frontend/`:
     ```bash
     cd frontend
     npm install
     npm start
     ```
   - The frontend will run at [http://localhost:3000](http://localhost:3000).


## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.