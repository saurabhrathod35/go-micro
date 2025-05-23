# Go Microservices Starter

A modular and scalable microservices starter kit built with Go (Golang), designed to accelerate the development of distributed systems. Each service is self-contained, promoting independent development, testing, and deployment. The repository includes example services (authentication, broker, logger, mail) and a simple front-end, each in its own directory.

## 📁 Project Structure.....

```
go-micro/
├── authentication-service/   # Auth microservice (stub)
├── broker-service/           # Broker API microservice (implemented)
│   └── cmd/api/              # Entrypoint: main.go
├── front-end/                # Basic HTML front-end
│   └── cmd/web/              # Entrypoint: main.go
├── logger-service/           # Logger microservice (stub)
├── mail-service/             # Mail microservice (stub)
├── project/                  # Shared config/tools (Makefile, etc.)
```

## 🚀 Getting Started

### Prerequisites

- Go 1.16 or higher
- Git

### Clone the Repository

```bash
git clone https://github.com/saurabhrathod35/go-micro.git
cd go-micro
```

### Run Services Locally

Each microservice is independently runnable. Here's how to start the broker and front-end:

**Run Broker Service:**

```bash
cd broker-service/cmd/api
go run main.go
```

**Run Front-End:**

Open a new terminal:

```bash
cd front-end/cmd/web
go run main.go
```

## 🌐 Usage

### Broker Service

- **POST /** – Returns a JSON response:
  ```bash
  curl -X POST http://localhost:8081
  ```
  Response:
  ```json
  {"error":false,"message":"Hit the broker"}
  ```

- **GET /ping** – Health check endpoint:
  ```bash
  curl http://localhost:8081/ping
  ```

### Front-End

Visit [http://localhost:8082](http://localhost:8082) in a browser.

- Click the "Hit Broker" button to send a POST request to the broker.
- The response is displayed in the browser using JavaScript.

## 🛠 Services Overview

| Service              | Status      | Description                          |
|----------------------|-------------|--------------------------------------|
| Broker Service       | ✅ Implemented | Example API gateway                  |
| Authentication       | 🔧 Stub       | Placeholder for auth service         |
| Logger               | 🔧 Stub       | Central logging service              |
| Mail                 | 🔧 Stub       | Email/notification service           |
| Front-End            | ✅ Implemented | Simple UI to test broker             |

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repo
2. Create a new branch: `git checkout -b feature-name`
3. Make your changes
4. Run and test the services locally
5. Submit a pull request with a clear description

Please format code with `gofmt` and follow Go best practices.

---

Happy coding! 🚀
