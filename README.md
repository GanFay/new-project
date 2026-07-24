# 🚀 SplitCore — Telegram Expense Organizer (Monorepo)

[![SplitCore CI](https://github.com/GanFay/SplitCore/actions/workflows/ci.yml/badge.svg)](https://github.com/GanFay/SplitCore/actions/workflows/ci.yml)

**SplitCore** is a modern, high-performance Telegram bot ecosystem designed to automate shared expense tracking for groups of friends, travelers, or event organizers.

The project has evolved from a single monolithic bot into a **highly scalable, event-driven microservices architecture**, showcasing advanced patterns of synchronous (gRPC) and asynchronous (RabbitMQ) service communication.

## 🔥 The Core Idea
Users create "Funds" (events), invite friends via unique deep-links, and record their expenses. The bot automatically calculates the balance: who overpaid and who needs to settle their debt using a greedy matching algorithm to minimize transactions.

### 🎬 Demo
*(Note: For the best experience, view the `.mp4` video directly if the GIF is buffering)*

![SplitCore Demo](demo2.gif)

## 🏗 System Architecture & Flow

To ensure high performance and zero blockage of the main Telegram Bot long-polling loop, the system is decoupled into two primary services using the **Event Carried State Transfer (ECST)** and **Notification Gateway** patterns [3]:

1. **`split-core` (The Core Engine):** Handles the database transactions, state management, and user interaction.
2. **`split-notify` (The Notification Processor):** An asynchronous worker service that formats and dispatches personalized user alerts [3].

```text
  [ Telegram Bot ]
         │
         ▼  (User adds an expense)
┌─────────────────────────────────────────────────────────┐
│                     split-core                          │
│  1. Commits transaction to PostgreSQL database           │
│  2. Publishes "expense_created" JSON event to RabbitMQ  │──┐
└─────────────────────────────────────────────────────────┘  │
         ▲                                     (gRPC)        │ (Async event)
         │                                     [3, 5]        │ [3]
         │ (2. gRPC Call: Get targets)                       ▼
         │ (4. gRPC Call: SendTelegramNotification)   ┌──────────────┐
         └────────────────────────────────────────────│ split-notify │
                                                      └──────────────┘
```

### The Async Notification Lifecycle [3]:
1. A user logs a new expense in `split-core` [3].
2. `split-core` commits the transaction to **PostgreSQL** and immediately publishes an `expense_created` JSON event (including the fund name and creator's name) to **RabbitMQ** [3].
3. The background worker `split-notify` consumes the event [3].
4. `split-notify` makes a synchronous **gRPC** request to `split-core` to fetch the target users to notify (excluding the expense creator) [2, 3].
5. `split-notify` generates beautiful, personalized **HTML-formatted** notification templates [3, 6].
6. `split-notify` calls `split-core`'s gRPC server to dispatch the parsed messages back to Telegram asynchronously, keeping the main thread responsive [2, 3].

---

## 🛠 Tech Stack
* **Language:** Go (Golang) 1.26.2
* **Framework:** [telebot.v4](https://github.com/tucnak/telebot) (Telegram Bot API)
* **Message Broker:** [RabbitMQ 3.12+](rabbitmq:4.3.2-management) (Quorum Queues via `amqp091-go` for async decoupled events) [3]
* **RPC Framework:** google.golang.org/grpc(via Protocol Buffers v3 for high-performance synchronous calls) [2]
* **Databases:** PostgreSQL (via `pgx/v5` Connection Pool), Redis 8.6.2 (via `go-redis/v9` for persistent FSM state).
* **Monorepo Tools:** Go Workspaces (`go.work`), Docker Multi-stage builds, Docker Compose, GNU Make.
* **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate).

---

## 📂 Project Structure (Monorepo)

The repository uses a strict Monorepo layout with isolated Go modules and shared contract specifications:

```text
SplitCore/ (Repository Root)
├── .github/                # GitHub Actions CI/CD workflows
├── proto/                  # gRPC contract definitions (.proto files and generated Go code)
│   └── pb/
├── split-core/             # The main transaction engine & Telegram Bot (Clean Architecture)
│   ├── cmd/
│   ├── internal/
│   │   ├── config/         # Config loader (reads environment variables)
│   │   ├── delivery/       # Adapters: Telegram handler & gRPC Server
│   │   ├── repository/     # Data stores: Postgres, Redis & RabbitMQ Publisher
│   │   └── usecase/        # Core business workflow orchestrators
│   └── Dockerfile
├── split-notify/           # Asynchronous Notification worker service (Layered Architecture)
│   ├── cmd/
│   ├── internal/
│   │   ├── client/         # gRPC client for split-core
│   │   ├── config/         # Config loader (reads environment variables)
│   │   ├── consumer/       # RabbitMQ Event consumer
│   │   └── processor/      # Event orchestrator (unmarshals JSON & routes messages)
│   └── Dockerfile
├── docker-compose.yml      # Orchestrates Postgres, Redis, RabbitMQ, split-core, and split-notify
└── Makefile                # Unified development task automation
```

---

## 📍 Roadmap

### ✅ Completed Milestones
- [x] **Monorepo Migration:** Restructured monolithic project into modular, independent services.
- [x] **Event-Driven Architecture:** Decoupled transactional operations from notifications using **RabbitMQ**.
- [x] **gRPC Bidirectional Communication:** Implemented ultra-fast internal synchronous communication via Protocol Buffers.
- [x] **Quorum Queues:** Configured persistent, fault-tolerant RabbitMQ queues.
- [x] **Go Workspaces (`go.work`):** Seamless local multi-module development and dependency resolving.
- [x] **Advanced Docker Orchestration:** Multi-stage, cached Docker builds with shared context in Docker Compose.
- [x] **PostgreSQL & Redis Integration:** Schema migrations and Redis-based FSM state storage.
- [x] **Math & Settlement Engine:** Greedy algorithm minimizing transaction chains for shared debt settling.
- [x] **CI/CD Pipeline:** Automated tests and static analysis (`golangci-lint`) via GitHub Actions.

### 🚀 Future Enhancements
- [ ] **Expense Management:** Ability to delete or edit logged mistakes.
- [ ] **Settle Debt feature:** "Mark as paid" logic to automatically adjust balances when someone returns the money.
- [ ] **Export to CSV:** Generate and download fund reports on the fly.
- [ ] **Multi-currency support.**

---

## ⚙️ Getting Started (Dev)

**Prerequisites:** Docker, Docker Compose, GNU Make, Protobuf Compiler (optional, for code gen).

1. Clone the repository:
   ```bash
   git clone https://github.com/GanFay/SplitCore.git
   ```
2. Set up environment variables. Copy the example file at the root and fill in your details (Bot token, DB credentials):
   ```bash
   cp .env.example .env
   ```
3. Initialize the Go Workspace locally (optional, for IDE support):
   ```bash
   go work init ./split-core ./split-notify ./proto
   ```
4. Start the infrastructure services (Postgres, Redis, RabbitMQ):
   ```bash
   make env-up
   ```
5. Run PostgreSQL database migrations:
   ```bash
   make migrate-up
   ```
6. Build and start all services in Docker Compose (including bot and notification processor):
   ```bash
   make run-services
   ```
7. To rebuild and hot-reload **only** the notification service after modifying its code:
   ```bash
   make run-notify
   ```