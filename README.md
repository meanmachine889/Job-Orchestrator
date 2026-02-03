# Distributed Job Orchestrator

A distributed job processing system built with Go and Next.js. This project implements a scalable architecture for managing background jobs with worker nodes, automatic retries, and real-time monitoring.

![Architecture](https://img.shields.io/badge/Architecture-Microservices-blue)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Next.js](https://img.shields.io/badge/Next.js-15-black?logo=next.js)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis)

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [API Reference](#api-reference)
- [Project Structure](#project-structure)

---

## Overview

This distributed job orchestrator allows you to:

1. **Submit jobs** via REST API or the dashboard UI
2. **Queue jobs** in Redis for reliable processing
3. **Process jobs** across multiple worker nodes
4. **Monitor** job status and worker health in real-time
5. **Handle failures** with automatic retries and dead-letter tracking

---

## Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   Dashboard     │────▶│   Orchestrator   │◀────│    Workers      │
│   (Next.js)     │     │   (Go HTTP API)  │     │   (Go Nodes)    │
└─────────────────┘     └────────┬─────────┘     └────────┬────────┘
                                 │                        │
                    ┌────────────┼────────────┐           │
                    ▼            ▼            ▼           │
              ┌──────────┐ ┌──────────┐                   │
              │PostgreSQL│ │  Redis   │◀──────────────────┘
              │  (Jobs)  │ │ (Queue)  │
              └──────────┘ └──────────┘
```

### Components

| Component | Description |
|-----------|-------------|
| **Orchestrator** | Central API server that manages jobs, workers, and job assignments |
| **Worker** | Stateless worker nodes that fetch and execute jobs from the queue |
| **Dashboard** | Real-time web UI to monitor jobs and workers |
| **PostgreSQL** | Persistent storage for jobs, workers, and job attempts |
| **Redis** | Message queue for job distribution using LPUSH/BRPOP |

---

## Features

### ✅ Job Management
- **Create jobs** with custom type, payload, and configuration
- **List jobs** with pagination support
- **Job details** view with full execution history
- **Job statuses**: `PENDING`, `RUNNING`, `SUCCESS`, `FAILED`, `RETRYING`, `DEAD`

### ✅ Automatic Retry System
- Configurable **max retries** per job
- Exponential backoff on failures
- Jobs marked as `DEAD` after exhausting all retries
- Error tracking for each attempt

### ✅ Worker Management
- **Dynamic worker registration** - workers self-register on startup
- **Heartbeat monitoring** - workers send heartbeat every 5 seconds
- **Automatic offline detection** - workers marked offline after 15 seconds of inactivity
- Worker status tracking (`ONLINE` / `OFFLINE`)

### ✅ Real-time Dashboard
- **Jobs page** - View all jobs with status badges
- **Workers page** - Monitor active workers
- **Job details** - Inspect payload, errors, and retry count
- **Create job dialog** - Submit new jobs from the UI
- Auto-refresh every 2-3 seconds

### ✅ Redis-based Job Queue
- Jobs pushed to Redis list on creation
- Workers block-wait on Redis for new jobs (BRPOP)
- Ensures reliable job distribution across workers

---

## Tech Stack

### Backend
- **Go** - Orchestrator and Worker services
- **PostgreSQL 16** - Primary database for jobs and workers
- **Redis 7** - Job queue
- **golang-migrate** - Database migrations

### Frontend
- **Next.js 15** - React framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **shadcn/ui** - UI component library

### Infrastructure
- **Docker Compose** - Local development environment

---

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- pnpm (or npm/yarn)

### 1. Start Infrastructure

```bash
# Start PostgreSQL and Redis
docker-compose up -d
```

This starts:
- PostgreSQL on `localhost:5432`
- Redis on `localhost:6379`

### 2. Setup Environment

Create a `.env` file in the project root:

```env
# Database
DATABASE_URL=postgres://orchestrator:orchestrator@localhost:5432/orchestrator?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# Orchestrator (for workers)
ORCHESTRATOR_URL=http://localhost:8080
```

### 3. Run Database Migrations

```bash
cd backend/orchestrator

# Install golang-migrate if needed
# brew install golang-migrate (macOS)
# Or download from https://github.com/golang-migrate/migrate

migrate -path migrations -database "$DATABASE_URL" up
```

### 4. Start the Orchestrator

```bash
cd backend/orchestrator
go run cmd/main.go
```

The orchestrator will start on `http://localhost:8080`

### 5. Start Worker(s)

Open a new terminal:

```bash
cd backend/worker
go run cmd/main.go
```

You can start multiple workers in separate terminals for distributed processing.

### 6. Start the Dashboard

```bash
cd frontend/dashboard
pnpm install
pnpm dev
```

Open `http://localhost:3000` to view the dashboard.

---

## API Reference

### Jobs

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/jobs` | Create a new job |
| `GET` | `/jobs` | List all jobs (with `?limit=` and `?offset=`) |
| `GET` | `/jobs/{id}` | Get job details by ID |
| `POST` | `/jobs/next` | Assign next pending job to a worker |
| `POST` | `/jobs/report` | Report job result (SUCCESS/FAILED) |

#### Create Job Request
```json
{
  "type": "email",
  "payload": { "to": "user@example.com", "subject": "Hello" },
  "max_retries": 3,
  "timeout_seconds": 30
}
```

#### Job Response
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "type": "email",
  "status": "PENDING",
  "retry_count": 0,
  "max_retries": 3,
  "timeout_seconds": 30,
  "created_at": "2026-02-03T10:00:00Z",
  "updated_at": "2026-02-03T10:00:00Z"
}
```

### Workers

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/workers/register` | Register a new worker |
| `POST` | `/workers/heartbeat` | Send worker heartbeat |
| `GET` | `/workers` | List all workers |

---

## Project Structure

```
distributed-job/
├── docker-compose.yml          # PostgreSQL & Redis
├── go.work                     # Go workspace
│
├── backend/
│   ├── orchestrator/           # Central API server
│   │   ├── cmd/main.go         # Entry point
│   │   ├── internal/
│   │   │   ├── api/            # HTTP handlers & routes
│   │   │   ├── queue/          # Redis queue operations
│   │   │   ├── scheduler/      # Worker health monitor
│   │   │   └── store/          # Database operations
│   │   └── migrations/         # SQL migrations
│   │
│   ├── worker/                 # Worker node
│   │   ├── cmd/main.go         # Entry point
│   │   └── internal/
│   │       ├── executor/       # Job execution logic
│   │       ├── orchestrator/   # API client
│   │       └── redis/          # Redis client
│   │
│   └── shared/                 # Shared models
│       └── job/model.go
│
└── frontend/
    └── dashboard/              # Next.js dashboard
        └── src/
            ├── app/            # Pages (jobs, workers)
            ├── components/     # UI components
            └── types/          # TypeScript types
```

---

## Job Types

The worker currently supports these job types:

| Type | Behavior |
|------|----------|
| `email` | Simulates email sending (2s delay, always succeeds) |
| `fail` | Simulates failure (1s delay, always fails) |
| `*` (default) | Generic job (1s delay, always succeeds) |

Extend the executor in `backend/worker/internal/executor/executor.go` to add custom job types.

---

## Database Schema

### Jobs Table
| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `type` | TEXT | Job type identifier |
| `payload` | JSONB | Job data/parameters |
| `status` | ENUM | PENDING, RUNNING, SUCCESS, FAILED, RETRYING, DEAD |
| `retry_count` | INT | Current retry attempt |
| `max_retries` | INT | Maximum retry attempts |
| `timeout_seconds` | INT | Job timeout |
| `worker_id` | UUID | Assigned worker (nullable) |
| `error` | TEXT | Error message (nullable) |
| `created_at` | TIMESTAMPTZ | Creation timestamp |
| `updated_at` | TIMESTAMPTZ | Last update timestamp |

### Workers Table
| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `hostname` | TEXT | Worker hostname |
| `status` | TEXT | ONLINE / OFFLINE |
| `last_heartbeat` | TIMESTAMPTZ | Last heartbeat time |

### Job Attempts Table
| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `job_id` | UUID | Foreign key to jobs |
| `attempt_number` | INT | Attempt number |
| `status` | ENUM | Attempt status |
| `error` | TEXT | Error message |
| `started_at` | TIMESTAMPTZ | Start time |
| `finished_at` | TIMESTAMPTZ | End time |

---

## License

MIT

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
