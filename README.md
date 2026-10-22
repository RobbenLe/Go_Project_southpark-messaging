# Go_Project_southpark-messaging
# South Park Messaging 📮

A small distributed system where South Park characters send each other messages **asynchronously**.

- **Go API** (Hexagonal Architecture) receives messages over HTTP and publishes them to RabbitMQ
- **RabbitMQ** stores them in a queue until someone reads them
- **Python consumer** reads them from the queue and prints them
- **Python spammer** sends a random South Park message every few seconds

## Architecture

```
[ spammer ]  (Python, every 3s)
     │ POST /messages {"author","body"}
     ▼
[ go-api ]  (Go, Hexagonal)
     │ validate → publish
     │ exchange "southpark" (topic)
     │ routing key "message.<author>"
     ▼
[ rabbitmq ]
     │ binding "message.#"
     ▼
[ queue: southpark_messages ]
     │
     ▼
[ consumer ]  (Python)
     prints: [message.cartman] Cartman: "Respect my authoritah!"
```

### Inside the Go API (Ports & Adapters)

```
HTTP adapter ──> (MessageSender port) ──> App service ──> (MessagePublisher port) ──> RabbitMQ adapter
 driving                                  validate +                                  driven
 adapter                                  publish                                     adapter
```

The core (`domain` + `app`) does not know about HTTP or RabbitMQ. It only knows the ports (interfaces).
Adapters plug into the ports from outside, and `cmd/api/main.go` wires everything together.

## Run it

Requirements: **Docker Desktop**.

```bash
docker compose up --build
```

This starts 4 containers: `rabbitmq`, `go-api`, `consumer` and `spammer`.
Within a few seconds the consumer starts printing the spammer's messages.

| What | URL |
|---|---|
| Go API | http://localhost:8080/messages |
| RabbitMQ management UI | http://localhost:15672 (user `southpark`, password `southpark`) |

Stop with `Ctrl+C`, then clean up:

```bash
docker compose down
```

## Try it yourself

**curl (macOS / Linux / Git Bash):**

```bash
curl -i -X POST http://localhost:8080/messages \
  -H "Content-Type: application/json" \
  -d '{"author":"Cartman","body":"Respect my authoritah!"}'
```

**PowerShell (Windows):**

```powershell
Invoke-RestMethod -Method Post -Uri http://localhost:8080/messages -ContentType "application/json" -Body '{"author":"Cartman","body":"Respect my authoritah!"}'
```

### Responses

| Request | Status | Body |
|---|---|---|
| Valid message | `202 Accepted` | `{"status":"queued"}` |
| Missing `author` | `400 Bad Request` | `{"error":"author is required"}` |
| Missing `body` | `400 Bad Request` | `{"error":"body is required"}` |
| Invalid JSON / unknown field | `400 Bad Request` | `{"error":"invalid JSON: ..."}` |
| RabbitMQ unavailable | `500 Internal Server Error` | `{"error":"could not send message"}` |

## Project structure

```
.
├── docker-compose.yml
├── go-api/
│   ├── Dockerfile               multi-stage build
│   ├── cmd/api/main.go          composition root (wiring + config)
│   └── internal/
│       ├── domain/              Message + validation rules
│       ├── ports/               MessageSender, MessagePublisher (interfaces)
│       ├── app/                 SendMessage use case
│       └── adapters/
│           ├── httpadapter/     POST /messages (driving adapter)
│           └── rabbitmq/        topic publisher (driven adapter)
├── consumer/                    Python consumer (pika)
└── spammer/                     Python timed message sender (requests)
```

| Layer | Depends on | Responsibility |
|---|---|---|
| `domain` | nothing | What a message is, and when it is valid |
| `ports` | domain | Interfaces between the core and the outside world |
| `app` | domain, ports | The "send a message" use case |
| `adapters` | everything | HTTP and RabbitMQ technology |
| `cmd/api` | everything | Builds the adapters and injects them |

Dependencies always point **inward**, towards the domain.

## Design decisions

- **Hexagonal Architecture:** business rules are isolated from technology. RabbitMQ could be replaced (e.g. by Kafka) by writing one new adapter, without touching `domain` or `app`. It also makes the core testable with a fake publisher.
- **Two ports:** `MessageSender` (driving: HTTP → app) and `MessagePublisher` (driven: app → RabbitMQ).
- **Topic exchange:** messages are published to the `southpark` topic exchange with routing key `message.<author>` (e.g. `message.cartman`). The queue is bound with `message.#`. New queues can subscribe to specific characters (e.g. `message.cartman`) without changing the API.
- **`202 Accepted`:** the API only queues the message; it is processed later by the consumer, so `200 OK` ("done") would be inaccurate.
- **Validation in the domain:** missing fields are a business rule, so they live in `domain.Validate()`. Malformed JSON is a transport problem, so it is handled in the HTTP adapter.
- **Reliability:**
  - durable exchange and queue, persistent messages (they survive a RabbitMQ restart)
  - manual `ack` after printing (a crashed consumer doesn't lose messages)
  - invalid messages are rejected with `nack(requeue=False)` to avoid poison-message loops
  - `prefetch_count=1` (the consumer handles one message at a time)
- **Startup order:** RabbitMQ has a Docker healthcheck; `go-api` and `consumer` only start once it is healthy. The consumer and spammer also retry their connections.
- **Configuration via environment variables** (`RABBITMQ_URL`, `API_URL`, ...): the same code runs locally (`localhost`) and in Docker (service names).

## Tests

The Go core is covered by unit tests (no RabbitMQ needed; the publisher is faked):

```bash
cd go-api
go test ./... -v
```

- `domain`: table-driven tests for `Validate()`
- `app`: `SendMessage` with a fake publisher (valid, invalid → not published, publisher failure)

## Bonus

Sometimes Kenny tries to send a message without a body. The API rejects it with `400`.
*Oh my God, they killed Kenny's message!*