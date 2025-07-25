# WASAText - Real-Time Messaging Application

[![Go](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Vue.js](https://img.shields.io/badge/Vue.js-3.x-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
[![SQLite](https://img.shields.io/badge/SQLite-003B57?style=flat&logo=sqlite)](https://www.sqlite.org/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker&logoColor=white)](https://www.docker.com/)

A modern, real-time messaging application built as part of the **Web and Software Architecture** course at Sapienza University of Rome. Built upon the "Fantastic Coffee (Decaffeinated)" template, WASAText extends the base framework to provide a comprehensive chat experience with support for individual and group conversations, media sharing, message reactions, and real-time updates.

## 🌟 Features

### Core Messaging
- **Real-time messaging** with automatic polling for instant updates
- **Private conversations** between two users
- **Group chats**
- **Message replies**
- **Message forwarding** across conversations
- **Message reactions** with emoji support
- **Media sharing** (image upload and display)

### User Experience
- **Intuitive chat interface** with modern design
- **Conversation management** with preview and search
- **Group administration** (add/remove participants, change names)
- **User profiles** with customizable profile pictures
- **Real-time status indicators** (delivered, read)

### Technical Features
- **RESTful API** with comprehensive OpenAPI documentation
- **Real-time updates** without WebSocket complexity
- **Secure authentication** and authorization
- **Database persistence** with SQLite
- **Dockerized deployment** for easy setup
- **Scalable architecture** with clean separation of concerns

## 🛠 Technology Stack

### Backend
- **Go 1.19+** - High-performance HTTP server
- **httprouter** - Fast HTTP routing and middleware
- **SQLite** - Lightweight, embedded database
- **OpenAPI 3.0** - API documentation and validation

### Frontend
- **Vue.js 3** - Progressive JavaScript framework
- **Vite** - Fast build tool and development server
- **Bootstrap 5** - Responsive CSS framework
- **Axios** - HTTP client for API communication

### Infrastructure
- **Docker** - Containerization and deployment
- **Docker Compose** - Multi-container orchestration
- **Go Vendoring** - Dependency management

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Node.js 20+ (for development)
- Go 1.19+ (for development)

### Using Docker (Recommended)
```bash
# Clone the repository
git clone https://github.com/Reewd/WASAproject.git
cd WASAproject

# Start the application
docker-compose up --build

# Access the application
# Frontend: http://localhost:80
# Backend API: http://localhost:3000
```

### Development Setup

#### Backend Development
```bash
# Run the Go backend
go run ./cmd/webapi/

# The API will be available at http://localhost:3000
```

#### Frontend Development
```bash
# Start development server
./open-node.sh
# Inside the container:
yarn run dev

# The frontend will be available at http://localhost:5173
```

### Production Build
```bash
# Build for production
./open-node.sh
# Inside the container:
yarn run build-prod

# Build Go binary with embedded frontend
# Inside the container:
yarn run build-prod

# Build Go binary with embedded frontend
go build -tags webui ./cmd/webapi/
```

## 📁 Project Structure

```
├── cmd/                    # Executable applications
│   ├── healthcheck/       # Health check daemon
│   └── webapi/           # Main web API server
├── doc/                   # API documentation
│   └── api.yaml          # OpenAPI 3.0 specification
├── service/              # Core application logic
│   ├── api/             # HTTP handlers and routing
│   ├── database/        # Database layer and models
│   └── globaltime/      # Time utilities for testing
├── webui/               # Vue.js frontend application
│   ├── src/            # Source code
│   │   ├── components/ # Reusable Vue components
│   │   ├── views/      # Page components
│   │   ├── services/   # API client services
│   │   └── modals/     # Modal dialogs
│   └── public/         # Static assets
├── vendor/             # Go dependencies (vendored)
├── docker-compose.yml  # Multi-container setup
└── Dockerfile.*       # Container definitions
```

## 🔌 API Documentation

The application provides a comprehensive RESTful API documented with OpenAPI 3.0. All endpoints (except `/session` and `/upload`) require authentication via Bearer token using the user ID.

**Base URL**: `http://localhost:3000`

### Authentication & Session
- `POST /session` - Login or create user account
- Authentication uses user ID as Bearer token: `Authorization: {userId}`

### User Management
- `GET /users` - List all users in the system
- `PUT /me/username` - Update current user's username
- `PUT /me/photo` - Update current user's profile photo

### Image Management
- `POST /upload` - Upload image file (multipart/form-data)
- `GET /uploads/{filepath}` - Serve uploaded images (static file serving)

### Conversations
- `GET /conversations` - List user's conversations
- `POST /conversations` - Create new conversation (private or group)
- `GET /conversations/{conversationId}` - Get conversation details with messages

### Group Management
- `PUT /conversations/{conversationId}/name` - Update group name
- `PUT /conversations/{conversationId}/photo` - Update group photo
- `POST /conversations/{conversationId}/participants` - Add participants to group
- `DELETE /conversations/{conversationId}/participants` - Leave group

### Messages
- `POST /conversations/{conversationId}/messages` - Send message (text, photo, or both)
- `DELETE /conversations/{conversationId}/messages/{messageId}` - Delete message
- `POST /conversations/{conversationId}/forwarded_messages` - Forward message to conversation

### Message Reactions
- `POST /conversations/{conversationId}/messages/{messageId}/reactions` - Add/update reaction
- `DELETE /conversations/{conversationId}/messages/{messageId}/reactions` - Remove reaction

### System
- `GET /liveness` - Health check endpoint

## 🏗 Architecture

### Backend Architecture
The Go backend follows a clean architecture pattern:

- **API Layer** (`service/api/`) - HTTP handlers, routing, and middleware
- **Database Layer** (`service/database/`) - Data access and persistence
- **DTOs** (`service/api/dto/`) - Data transfer objects
- **Validation** (`service/api/constraints/`) - Input validation rules

### Frontend Architecture
The Vue.js frontend uses a component-based architecture:

- **Components** - Reusable UI components (Chat, Message, ConversationPreview)
- **Views** - Page-level components
- **Composables** - Shared reactive logic (authentication, validation)
- **Services** - API communication layer

### Database Schema
SQLite database with optimized schema for:
- Users and authentication
- Conversations (private and group)
- Messages
- Reactions and media attachments
- Read receipts and delivery status

## 🔧 Development

### Code Organization
- **Go Modules** - Dependency management with vendoring
- **Vue 3 Composition API** - Modern reactive patterns
- **TypeScript-style JSDoc** - Enhanced developer experience
- **ESLint/Prettier** - Code formatting and linting

### Development Workflow
```bash
# Backend development
go run ./cmd/webapi/
go test ./...
go test ./...

# Frontend development
# Frontend development
./open-node.sh
yarn run dev

# Database migrations
# Automatic schema updates on startup

# Dependency management
go mod tidy && go mod vendor
yarn install --immutable

# Database migrations
# Automatic schema updates on startup

# Dependency management
go mod tidy && go mod vendor
yarn install --immutable
```

## 🧪 Testing & Quality

### Features Tested
- **API endpoints** with comprehensive request/response validation
- **Database operations** with transaction safety
- **Authentication and authorization** flows
- **Real-time updates** and polling mechanisms
- **Media upload and serving** functionality

## 🎯 Key Learning Outcomes

This project demonstrates proficiency in:

- **Full-stack web development** with modern technologies
- **RESTful API design** and implementation
- **Real-time application architecture** without WebSockets
- **Database design** and optimization
- **Docker containerization** and deployment
- **Frontend state management** and reactive programming
- **Security best practices** and authentication
- **Code organization** and clean architecture principles

## 📚 Course Context

Built for the **Web and Software Architecture** course at Sapienza University of Rome, this project extends the provided "Fantastic Coffee (Decaffeinated)" template to showcase:
- Understanding of web application architecture patterns
- Implementation of modern development practices
- Ability to integrate multiple technologies effectively
- Knowledge of deployment and containerization strategies
- Full-stack development skills from template to complete application

## 🚀 Deployment

### Production Deployment
```bash
# Build production frontend
./open-node.sh
yarn run build-prod
exit

# Build Go binary with embedded frontend
go build -tags webui ./cmd/webapi/

# Run production server
./webapi
```

### Docker Deployment
```bash
# Production deployment
docker-compose up -d

# Health check
curl http://localhost:3000/liveness
```

## 📄 License

This project is based on the "Fantastic Coffee (Decaffeinated)" template provided by the course instructors under the MIT License. 

**Base Template**: [sapienzaapps/fantastic-coffee-decaffeinated](https://github.com/sapienzaapps/fantastic-coffee-decaffeinated)

**Original Copyright**: Copyright (c) 2022 Enrico Bassetti

The original template is licensed under the MIT License. All modifications and extensions for the WASAText messaging application were developed as part of academic coursework at Sapienza University of Rome.

## 🙏 Acknowledgments

- **Enrico Bassetti** - Original "Fantastic Coffee (Decaffeinated)" template
- **Emanuale Panizzi** - Course instructor
- **Sapienza University of Rome** - Web and Software Architecture course
- **Course instructors** for providing the project framework and guidance
- **Go community** for excellent tooling and libraries
- **Vue.js team** for the reactive framework

---

**Note**: This is an educational project developed for the Web and Software Architecture course, built upon the "Fantastic Coffee (Decaffeinated)" template provided by course instructors. While fully functional, it's designed for learning purposes and academic evaluation.
