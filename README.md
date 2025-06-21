# Service Exchange Platform
## 💡 Overview

The Service Exchange Platform is a microservices-based marketplace designed to connect users who need tasks done with service providers (taskers) who can complete them. It manages the entire task lifecycle, from creation and bidding to payment and reviews.

Built upon **Domain-Driven Design (DDD)** principles and an **Event-Driven Architecture**, the platform emphasizes flexibility, scalability, and maintainability.

## ✨ Key Features

  * **User & Profile Management:** Registration, authentication, identity verification, profile updates.
  * **Task Management:** Creation, categorization, search, filtering, and status management of tasks.
  * **Bidding & Negotiation:** Handling bids, tasker selection, and communication.
  * **Payment System:** Secure financial transactions, escrow management, fee handling.
  * **Review & Rating System:** Reputation building and feedback mechanisms.
  * **Notifications & Communication:** Multi-channel alerts (in-app, push, email, SMS).
  * **Location Services:** Geospatial data management and proximity matching.
  * **Support & Dispute Resolution:** Ticketing system and conflict resolution processes.

## 🚀 Architectural Snapshot

This platform leverages a loosely coupled microservices architecture, featuring:

  * **API Gateway:** Single entry point for all client requests (e.g., Kong/Nginx).
  * **Service Mesh:** For inter-service communication management and security (e.g., Istio).
  * **Event Bus:** Apache Kafka for asynchronous communication between microservices.
  * **Database per Service:** Each microservice owns its dedicated data store (PostgreSQL, MongoDB, PostGIS, Redis).
  * **Containerization:** Docker for packaging applications and Kubernetes for orchestration.

### Core Technologies

| Category                    | Technologies Used/Recommended     |
| :-------------------------- | :-------------------------------- |
| **API Gateway** | Kong, Nginx, AWS API Gateway      |
| **Backend Services** | Go, Node.js (TypeScript) |
| **Message Broker** | Apache Kafka, RabbitMQ, Cloud Pub/Sub |
| **Databases** | PostgreSQL, MongoDB, Redis, PostGIS |
| **Container Orchestration** | Kubernetes          |
| **Service Mesh** | Istio, Linkerd                    |
| **Monitoring** | Prometheus + Grafana, Loki   |
| **CI/CD** | GitHub Actions |

## 📚 Comprehensive Documentation

All detailed architectural documentation is hosted on our **[GitHub Wiki](https://github.com/pratchaya-maneechot/service-exchange/wiki/Service-Exchange-Platform)**, covering topics such as:

1.  Architecture Overview
2.  Domain Model & Bounded Contexts
3.  Service Catalog
4.  Integration Patterns
5.  Event-Driven Architecture
6.  Data Management
7.  Security Architecture
8.  Infrastructure & Deployment
9.  Monitoring & Observability
10. Development Guidelines
11. Troubleshooting Guide

## 🛠️ Getting Started (For Developers)

### Prerequisites

  * Docker
  * kubectl (if deploying to Kubernetes)
  * Nodejs (v22.15.0), Go (go1.24.3)

### Local Setup & Running (Example for a single service)

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/pratchaya-maneechot/service-exchange
    cd service-exchange
    ```
2.  **Run infra with Docker Compose:**
    ```bash
    docker compose up -d
    ```
    **Run a single service with Local:**
    ```bash
    # Example for User Service
    make users-generate
    make users-migrate-up # Make sure the database "users" already exists
    make users-dev
    ```
    **Run api-gateway with Local:**
    ```bash
    # Example for User Service
    cd apps/api-gateway
    cp .env.sample .env
    pnpm i
    pnpm start:dev
    ```
3.  **Access the API Gateway:**
    The API Gateway typically runs on `http://localhost:8080/graphql` or as per your configuration.

### Useful Commands

Find common Docker, Kubernetes, and Kafka commands in the **[Troubleshooting Guide](https://github.com/pratchaya-maneechot/service-exchange/wiki/Troubleshooting-Guide)** under section `C. Useful Commands`.

## 🤝 Contributing

We welcome contributions\! Please refer to our **[Development Guidelines](https://github.com/pratchaya-maneechot/service-exchange/wiki/Development-Guidelines)** and the **[Code Review Process](https://github.com/pratchaya-maneechot/service-exchange/wiki/Development-Guidelines#code-review-process)** for more information.

## 📄 License

This project is licensed under the [MIT License](https://www.google.com/search?q=LICENSE).