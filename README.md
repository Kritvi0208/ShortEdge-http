
# ShortEdge (net/http Version)

A full-stack, privacy-conscious URL shortener built using Go’s standard `net/http` package.  
Includes branded short links, real-time analytics, visibility controls, link expiry, and a custom frontend UI.

> 🔁 Looking for the GoFr-based version? [View ShortEdge-gofr](https://github.com/Kritvi0208/ShortEdge-gofr)

---

## 🚀 Features
- Custom short links with branded codes
- Public/private toggle per URL
- Real-time analytics (location, browser, device)
- Link expiry support
- Full CRUD REST API
- HTML/CSS/JS frontend for testing and interaction
- Prometheus metrics endpoint (`/metrics`)
- Health check endpoint (`/health`)

---

## 🛠️ Tech Stack
- Go (net/http)
- PostgreSQL
- ipwho.is (GeoIP)
- uasurfer (User-Agent parsing)
- HTML, CSS, JavaScript

---

## 📂 Project Structure
```
/main.go
/handlers
/service
/store
/templates (for frontend)
/static (CSS/JS)
````

## 🧪 Run Locally
```bash
go run main.go
````

Visit: `http://localhost:8080`

---

