# IP-Geolocation API

A production-ready REST API for high-performance IP tracking and geolocation built with Go. Designed for low-latency visit recording and automated lifecycle management — featuring public IP visit tracking, geolocation enrichment via ipgeolocation.io, SQLite persistence, and visit history management.

---

## Architecture Overview

```
Client
  │
  ▼
Gin HTTP Server
  │
  ├── Handlers → Services → Repositories
  │                             │
  │                         SQLite (Persistent IP visit & history store)
  │
  └── ipgeolocation.io (Third-party IP Geolocation API)
        ├── Resolves IP to country, city, timezone, ISP, etc.
        └── Enriches visit records on each tracked request
```

---

## Tech Stack

| Layer | Technology | Reason |
| --- | --- | --- |
| Language | Go | High performance, lightweight goroutines for worker processes |
| Framework | Gin | Fast HTTP routing and clean middleware chaining |
| Database | SQLite | Lightweight, file-based relational store for IP visit and history metadata |
| Geolocation | ipgeolocation.io | Third-party API for enriched IP geolocation data |

---

## Features

### IP Visit Tracking

* **Automatic IP Detection:** Captures the visitor's IP address on each request and records the visit with a timestamp.
* **Geolocation Enrichment:** Each tracked visit is enriched with geolocation data from ipgeolocation.io, including country, city, region, timezone, ISP, and more.

### Visit & History Management

* **Visit Records:** Stores individual visit events per IP address, queryable by ID, IP, or date range.
* **History Records:** Maintains a historical log of all IP interactions, queryable by ID, IP, or date.
* **Today Filtering:** Dedicated endpoints to fetch visits and histories recorded on the current day.

---

## API Endpoints

<details>
<summary><strong>Geo IP</strong> — 9 endpoints</summary>

<br>

<details>
<summary>&nbsp;&nbsp;&nbsp;&nbsp;<strong>Visits</strong> — 5 endpoints</summary>

<br>

| Method | Endpoint | Access | Notes |
| --- | --- | --- | --- |
| GET | `/api/v1/geo/ip` | Public | Track a visit & enrich with geolocation data |
| GET | `/api/v1/geo/ip/visits/:id` | Public | Get IP visit record by Visit ID |
| GET | `/api/v1/geo/ip/visits/:ip` | Public | Get IP visit record by IP Address |
| GET | `/api/v1/geo/ip/visits` | Public | Get all IP visit records |
| GET | `/api/v1/geo/ip/visits/today` | Public | Get all IP visit records from today |

</details>

<details>
<summary>&nbsp;&nbsp;&nbsp;&nbsp;<strong>Histories</strong> — 4 endpoints</summary>

<br>

| Method | Endpoint | Access | Notes |
| --- | --- | --- | --- |
| GET | `/api/v1/geo/ip/histories/:id` | Public | Get IP history record by History ID |
| GET | `/api/v1/geo/ip/histories/:ip` | Public | Get all IP history records by IP Address |
| GET | `/api/v1/geo/ip/histories` | Public | Get all IP history records |
| GET | `/api/v1/geo/ip/histories/today` | Public | Get all IP history records from today |

</details>

</details>

---

## Getting Started

### Prerequisites

* Go 1.22+

### Run Locally

```bash
# Clone the repository
git clone https://github.com/Tyomaaans/Api-Geolocation.git
cd ip-geolocation

# Copy environment variables
cp .env.example .env

# Start go
go run ./cmd/api

```

### Environment Variables

See `.env.example` for all required variables. Key configs:

```env
# App
APP_PORT=

# Database
DATABASE_URL=

# Geolocation
IPGEO_API_KEY=

```

> `IPGEO_API_KEY` is obtained from [ipgeolocation.io](https://ipgeolocation.io). The API response is used to enrich each visit record with geolocation metadata such as country, city, region, timezone, ISP, latitude, longitude, and more.

---

## Project Status

| Feature | Status | Notes |
| --- | --- | --- |
| IP Visit Tracking (`GET /geo/ip`) | ✅ Done | Captures and records visitor IP on request |
| Geolocation Enrichment | ✅ Done | Enriches records via ipgeolocation.io API |
| Get Visit by ID | ✅ Done | `/geo/ip/visits/:id` |
| Get Visit by IP | ✅ Done | `/geo/ip/visits/:ip` |
| Get All Visits | ✅ Done | `/geo/ip/visits` |
| Get Today's Visits | ✅ Done | `/geo/ip/visits/today` |
| Get History by ID | ✅ Done | `/geo/ip/histories/:id` |
| Get Histories by IP | ✅ Done | `/geo/ip/histories/:ip` |
| Get All Histories | ✅ Done | `/geo/ip/histories` |
| Get Today's Histories | ✅ Done | `/geo/ip/histories/today` |
| SQLite Persistence | ✅ Done | Lightweight file-based store for visits & histories |