# Go Microservices Starter

A modular and scalable microservices starter kit built with Go (Golang), designed to accelerate the development of distributed systems using clean architecture principles.

## 🚀 Overview

This repository serves as a foundational template for building microservices in Go. It incorporates essential components such as service discovery, logging, and containerization, facilitating rapid development and deployment of microservices.

## 🧱 Project Structure

The project is organized into the following directories:

- `authentication-service/` – Handles user authentication and authorization.
- `broker-service/` – Acts as an intermediary for inter-service communication.
- `logger-service/` – Manages centralized logging across services.
- `front-end/` – Provides a user interface for interacting with the microservices.
- `project/` – Contains shared configurations and utilities.

Each service is self-contained, promoting independent development, testing, and deployment.

## 🛠️ Features

- **Microservices Architecture**: Encourages separation of concerns and scalability.
- **Clean Codebase**: Follows Go best practices for maintainability.
- **Dockerized Services**: Simplifies deployment with Docker containers.
- **Centralized Logging**: Aggregates logs for easier monitoring and debugging.
- **Service Discovery**: Enables dynamic discovery of services within the ecosystem.

## 🧰 Technologies Used

- **Go (Golang)**: Primary language for service implementation.
- **Docker**: Containerization of services for consistent environments.
- **Makefile**: Automates build and deployment tasks.
- **Go Modules**: Manages dependencies for each service.

## 🚀 Getting Started

### Prerequisites

- Go 1.16 or higher
- Docker
- Make

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/saurabhrathod35/go-micro.git
   cd go-micro
