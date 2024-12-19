# Messenger Microservices Architecture

This project is a **Messenger Application** built using a **microservices architecture**. The system is designed to handle authentication, messaging, presence tracking, and gateway routing for real-time communication.

---

## **Table of Contents**
1. [Overview](#overview)
2. [Microservices](#microservices)
3. [Technology Stack](#technology-stack)
4. [Project Structure](#project-structure)
5. [Getting Started](#getting-started)
6. [API Endpoints](#api-endpoints)
7. [Environment Variables](#environment-variables)
8. [Contributing](#contributing)
9. [License](#license)

---

## **Overview**

This application consists of multiple microservices working together to provide a real-time messaging platform. Each service has a clearly defined role and communicates with other services via HTTP REST APIs.

---

## **Microservices**

1. **Auth Service**
    - Handles user authentication and token validation.
    - Provides APIs for user registration, login, and token validation.

2. **Message Service**
    - Manages channels and message history.
    - Provides APIs to create channels, post messages, and fetch message history.

3. **Presence Service**
    - Tracks user presence and activity within channels.
    - Broadcasts user join/leave events to the Gateway Service.

4. **Gateway Service**
    - Acts as a WebSocket entry point.
    - Routes messages to appropriate channels and services.

---

## **Technology Stack**

- **Programming Language**: Go (Golang)
- **Database**: SQLite (for lightweight storage)
- **Communication**: REST APIs and WebSockets
- **Orchestration**: Docker & Docker Compose
- **Frameworks/Libraries**:
    - `gorilla/mux` for routing
    - `gorilla/websocket` for WebSocket communication
    - `bcrypt` for password hashing
    - `JWT` for token-based authentication

---

## **Project Structure**




chat-application/
README.md
docker-compose.yml  # If using Docker Compose to orchestrate services
.env                # Environment variables for services (optional)

    gateway-service/
        cmd/
            gateway/                  # main entry point (main.go)
        internal/
            handlers/                 # WebSocket handlers, connection mgmt
            authclient/               # code for calling Auth service
            messageclient/            # code for calling Message service
            presenceclient/           # (optional) code for calling Presence service
        pkg/
            models/                   # shared data models (Message, User, etc.)
            utils/                    # utility functions (logging, config)
        go.mod
        go.sum

    message-service/
        cmd/
            message/                  # main entry point (main.go)
        internal/
            storage/                  # SQLite handling (queries, migrations)
            handlers/                 # HTTP handlers for messages/channels
            services/                 # application logic (store message, retrieve messages)
        pkg/
            models/                   # data models (Message, Channel)
            utils/                    # utility functions, config handling
        migrations/                   # SQL migration files if needed
        go.mod
        go.sum

    auth-service/
        cmd/
            auth/                     # main entry point (main.go)
        internal/
            storage/                  # SQLite user storage and queries
            handlers/                 # HTTP handlers for login, signup, token validation
            jwt/                      # JWT generation and validation logic
        pkg/
            models/                   # data models (User, TokenClaims)
            utils/                    # utility functions, password hashing, config
        migrations/                   # SQL migration files for user db
        go.mod
        go.sum

    presence-service/ (optional)
        cmd/
            presence/                 # main entry point (main.go)
        internal/
            handlers/                 # possibly REST endpoints to query presence
            memory/                   # in-memory data structures for presence tracking
            broadcaster/              # logic to notify Gateway service of changes
        pkg/
            models/                   # presence models (OnlineUser)
            utils/                    # utility functions
        go.mod
        go.sum

    frontend/
        build/                        # compiled build of your chosen front-end framework (React, Vue, etc.)
        src/                          # front-end source code
        public/                       # static files
        package.json
        package-lock.json




---

## **Getting Started**

### Prerequisites

- [Go](https://golang.org/dl/) (1.17 or higher)
- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- [Postman](https://www.postman.com/) (optional, for API testing)

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/genryusaishigikuni/messenger.git
   cd messenger








2.Set up environment variables:
Create a .env file for each service


3.Run services with Docker Compose:

docker-compose up --build





4.Access services:

Auth Service: http://localhost:8082
Message Service: http://localhost:8081
Presence Service: http://localhost:8083
Gateway Service: http://localhost:8080








API Endpoints
Auth Service (http://localhost:8082)
POST /api/auth/register: Register a new user.
POST /api/auth/login: Login and receive a JWT token.
GET /api/auth/validate: Validate a token.


Message Service (http://localhost:8081)
GET /api/channels: Retrieve all channels.
POST /api/channels: Create a new channel.
GET /api/messages/history: Get message history for a channel.
POST /api/messages: Post a message to a channel.


Presence Service (http://localhost:8083)
GET /api/presence: Get all online users.
POST /api/presence/join: Mark a user as online in a channel.
POST /api/presence/leave: Mark a user as offline.


Gateway Service (http://localhost:8080)
GET /ws: WebSocket endpoint for real-time messaging.
