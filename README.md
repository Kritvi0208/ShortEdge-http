# 🔗 ShortEdge – A URL Shortener (net/http Version)

**ShortEdge** is a full-stack, privacy-conscious URL shortener built using Go’s `net/http` package.  
It supports branded links, real-time analytics, visibility toggles, link expiry, and Prometheus metrics — all powered by a clean HTML/CSS frontend and PostgreSQL backend.



> 🔧 Fully working MVP version using `net/http`  
> ⚙️ Production-ready, extensible, and built with backend performance + observability in mind
> 
> 🔁 Looking for the GoFr-native version? Check out [ShortEdge-gofr](https://github.com/Kritvi0208/ShortEdge)


---

## 🧩 Core Features

| Feature                  | Description                                                                 |
|-------------------------|-----------------------------------------------------------------------------|
| 🔗 Branded Short Links   | Supports custom short codes (e.g., `/r/my-event`)                           |
| 👁️ Public/Private Toggle  | Control analytics visibility per short link                                 |
| ⏳ Link Expiry Support    | Set expiry time after which the link auto-deactivates                      |
| 📊 Real-Time Analytics   | Logs IP, browser, device, and location for every redirect                   |
| ⚙️ RESTful CRUD API      | Endpoints for create, read, update, delete operations                       |
| 🧠 Device + Geo Parsing   | Integrated `ipwho.is` + `uasurfer` based visit logging                     |
| 📈 Prometheus Metrics    | Operational metrics for observability and monitoring                       |
| 🌐 HTML/CSS Frontend     | Lightweight UI built using Go’s `html/template` + static CSS               |
| 💡 Health Check Route    | `/health` confirms DB + server uptime                                      |

---


## 🛠️ System Architecture

```text
┌────────────────────────────┐
│ 🌍 Client (Web / API)      │
│ ────────────────────────── │
│ • HTML Web UI              │
│ • Postman / API Clients    │
└─────────────┬──────────────┘
              ▼
┌────────────────────────────┐
│ net/http Router            │
│ • Matches routes           │
│ • Directs to handler funcs │
└─────────────┬──────────────┘
              ▼
┌────────────────────────────┐
│ Handler Layer              │
│ • Parses input             │
│ • Validates + responds     │
└─────────────┬──────────────┘
              ▼
┌────────────────────────────┐
│ Logic Layer (Inline)       │
│ • Code generation          │
│ • Expiry, visibility check │
└─────────────┬──────────────┘
              ▼
┌────────────────────────────┐
│ PostgreSQL Storage         │
│ • Stores URL + metadata    │
│ • Logs visit info          │
└─────────────┬──────────────┘
              ▼
┌────────────────────────────┐
│ Analytics Hooks            │
│ • ipwho.is → Country       │
│ • uasurfer → Browser/Device│
└────────────────────────────┘
````
---
## 🗂️ Project Structure

```
ShortEdge/
├── frontend/                # Static frontend files (HTML + JS + CSS)
│   ├── all.html             # Lists all shortened URLs
│   ├── analytics.html       # Displays analytics for a specific short code
│   ├── delete.html          # UI for deleting a short URL
│   ├── health.html          # Health check status page
│   ├── index.html           # Main page for shortening URLs
│   ├── update.html          # UI to update a short URL
│   ├── script.js            # Shared JavaScript functionality
│   └── style.css            # Common styling
│
├── internal/
│   └── model/               # Domain data models
│       ├── url.go           # URL struct (short → long)
│       └── visit.go         # Visit struct (analytics: IP, country, device)
│
├── repository/              # Data access layer (PostgreSQL queries)
│   ├── db.go                # DB connection setup
│   ├── shortcode.go         # Short code generator logic
│   ├── url.repository.go    # CRUD ops for URL model
│   └── visit.repository.go  # CRUD ops for Visit logs
│
├── utils/                   # Helper utilities
│   └── parse.go             # User-Agent parsing logic
│
├── migration/
│   └── 001_create_tables.sql # SQL script to create required tables
│
├── .env                     # Environment variables (DB_URL, etc.)
├── go.mod                   # Go module definition
├── go.sum                   # Go module checksums
├── main.go                  # Application entrypoint (can be moved to cmd/main.go)
└── README.md                # Project overview and documentation
```
---

## 🔌 API Endpoints

| Method | Endpoint            | Description                                |
| ------ | ------------------- | ------------------------------------------ |
| POST   | `/shorten`          | Create a new branded or random short link  |
| GET    | `/r/{code}`         | Redirects and logs analytics               |
| GET    | `/analytics/{code}` | Returns all visits for the given short URL |
| GET    | `/all`              | List all non-expired short links           |
| PUT    | `/update/{code}`    | Update target URL or visibility            |
| DELETE | `/delete/{code}`    | Delete a short link                        |
| GET    | `/metrics`          | Exposes Prometheus metrics                 |
| GET    | `/health`           | Confirms DB + backend health               |

---
## 🧪 Tech Stack

| Layer        | Stack / Tool                                   |
| ------------ | ---------------------------------------------- |
| Backend      | Go (`net/http`)                                |
| Frontend     | HTML + CSS                                     |
| Database     | PostgreSQL                                     |
| Device Info  | [`uasurfer`](https://github.com/avct/uasurfer) |
| Geo Location | [`ipwho.is`](https://ipwho.is)                 |
| Monitoring   | [Prometheus](https://prometheus.io)            |

---

## 📸 Demo Screenshots

### 1. 🌐 Frontend user interface (`net/http`)
A minimal, responsive HTML/CSS interface for submitting long URLs, choosing custom short codes, toggling visibility, and receiving branded short links.

![Frontend UI](assets/ui-home.png)

---

### 2. 📊 `GET /analytics/{code}`
Returns rich, real-time analytics per short link — including:
- Visitor country (via [ipwho.is](https://ipwho.is))
- Device and browser (via [uasurfer](https://github.com/avct/uasurfer))
- Timestamp

![Analytics Endpoint](assets/get-analytics.png)

---

### 3. 📄 `GET /all`  
Lists all shortened links (public/private) with long URL mapping.

![All Links](assets/get-all.png)

---

### 4. 📈 Prometheus `/metrics`  
Live metrics exported at `/metrics`, ready for Grafana or alerting dashboards.

![Prometheus Metrics](assets/metrics-page.png)

---

## 🚀 Getting Started

### 🔧 Requirements

* Go 1.20+
* PostgreSQL v13+
* Internet access (for geo/device APIs)

### 🛠 Setup

1. **Clone the repo**

```bash
git clone https://github.com/Kritvi0208/ShortEdge-http
cd ShortEdge-http
```

2. **Configure PostgreSQL**
   Update your connection string inside `main.go`:

```go
dsn := "postgres://username:password@localhost:5432/shortedge?sslmode=disable"
```

3. **Run the app**

```bash
go run main.go
```

4. **Open in browser**

* UI: `http://localhost:8080`
* Metrics: `http://localhost:8080/metrics`

---

## 🏗 Real-World Use Cases

* 🔗 Custom short links for Google Forms, PDFs, feedback links
* 🔒 Private academic resource sharing
* 📈 Insight collection for link click-through rate
* 📊 Prometheus-ready analytics for observability dashboards

---

## Acknowledgements

* [ipwho.is](https://ipwho.is) — for free IP-to-country geolocation
* [uasurfer](https://github.com/avct/uasurfer) — for clean browser/device detection
* [Prometheus](https://prometheus.io) — for scalable metrics
---
