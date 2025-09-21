# Notes App Monorepo

This repository contains a full-stack notes application with a Go backend and a Next.js frontend.

## Project Structure

```
notes-backend/    # Go REST API server
notes-frontend/   # Next.js web client
docker-compose.yml
```

### Backend (`notes-backend`)

- Built with Go.
- Provides RESTful APIs for notes, users, categories, authentication, and image uploads.
- See [`notes-backend/README.md`](notes-backend/README.md) for backend-specific instructions.

### Frontend (`notes-frontend`)

- Built with Next.js.
- Uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) for optimized font loading.
- See [`notes-frontend/README.md`](notes-frontend/README.md) for frontend-specific instructions.

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Node.js (for frontend development)
- Go (for backend development)

### Development

Start both backend and frontend with Docker Compose:

```bash
docker-compose up --build
```

Or run each service manually:

#### Backend

```bash
cd notes-backend
go run cmd/main.go
```

#### Frontend

```bash
cd notes-frontend
npm install
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) for the frontend.

