# Snapp Chat

Snapp Chat is a real-time chat application written in Golang using NATS for broadcasting. The project consists of two main components:

1. **Server**: Handles the backend operations and WebSocket connections.
2. **Client**: A CLI-based client for interacting with the chat server.

---

## Features
- Create and join chatrooms.
- Real-time messaging using WebSocket.
- Command options for enhanced CLI user experience.

---

## Server

### Prerequisites
Ensure you have the following installed:
- [Git](https://git-scm.com/)
- [Docker and Docker Compose](https://docs.docker.com/)

### Installation
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd chatroom/server
   ```
2. Build and start the server:
   ```bash
   docker compose up -d --build
   ```

### Default Configuration
The default configuration file is as follows:
```json
{
  "nats" : {
    "host" : "snapp-chat-nats-server", // Container name of the NATS server
    "port" : 4222
  },
  "server" : {
    "host" : "", // Ensures Fiber listens on all network interfaces
    "port" : 8080
  }
}
```

### Endpoints

1. **Create Chatroom**
    - **Endpoint**: `POST /api/v1/chat/create`
    - **Request Body**:
      ```json
      {
        "creator": "your_username",
        "chatroomName": "chatroom_name"
      }
      ```  
    - **Response**:
      ```json
      {
        "chatroomId": "generated_chatroom_id",
        "userID": "creator_user_id"
      }
      ```

2. **Join Chatroom**
    - **Endpoint**: `POST /api/v1/chat/join`
    - **Request Body**:
      ```json
      {
        "chatRoomId": "chatroom_id",
        "name": "your_username"
      }
      ```  
    - **Response**:
      ```json
      {
        "userId": "joined_user_id",
        "Name": "your_username"
      }
      ```

3. **WebSocket Connection**
    - **Endpoint**: `chatroom/:chatId?userId=id`
        - `chatId`: The unique ID of the chatroom.
        - `userId`: The ID of the user connecting to the WebSocket.

### Error Format
All errors are returned as a JSON object:
```json
{
  "message": "error description"
}
```

---

## Client

### Features
The client is a Command Line Interface (CLI) application with several options for interacting with the server.

### Options
- **`-url`**: Specify the base URL of the server. (Default is **`localhost:8080`**)
- **`-u`**: Username for the client.
- **`-c`**: Create a chatroom.
- **`-i`**: Join a chatroom.

### Commands
- **`#exit`**: Exit the chatroom.
- **`#users`**: Display a list of users in the chatroom.

### Usage
1. Build the client:
   ```bash
   go build -o client ./client
   ```
2. Run the client with the following options:

    - **Create a Chatroom**:
      ```bash
      ./client -url localhost:8080 -u <username> -c <chatroom_name>
      ```

    - **Join a Chatroom**:
      ```bash
      ./client -url localhost:8080 -u <username> -i <chatroom_id>
      ```

---

## License
This project is licensed under the [MIT License](LICENSE).

Happy chatting! 🚀