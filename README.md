# 🧾 CMU Lifelong Education — Course Document Converter

A full-stack web application for extracting and displaying structured course information from `.docx` files.  
Developed for **Chiang Mai University School of Lifelong Education**.

---

## 🧩 Project Structure

| Directory | Description |
|------------|-------------|
| `back/` | Go (Fiber) backend API for document conversion |
| `front/` | Next.js frontend for uploading and displaying parsed data |
| `nginx/` | Reverse proxy configuration for production |
| `docker-compose.yml` | Containerized setup for both frontend and backend |

---

## 🚀 Getting Started (via Docker)

### 1️⃣ Prerequisites
- Docker & Docker Compose (v2+)
- Git

---

### 2️⃣ Clone the repository
```bash
git clone https://github.com/Ismailax/LE-Docx-Converter.git
cd LE-Docx-Converter
```

---

### 3️⃣ Configure environment variables

#### 🔧 Backend (`back/.env`)
```bash
PORT=
CORS_ALLOW_ORIGINS=
MAX_UPLOAD_MB=
```

#### 🔧 Frontend (`front/.env`)
```bash
NEXT_PUBLIC_APP_BASEPATH=
NEXT_PUBLIC_BACKEND_URL=
```

---

### 4️⃣ Run everything with Docker Compose

```bash
docker compose --env-file ./front/.env up -d --build
```

Docker will build and run 3 containers:

| Service | Internal Port | External Port | Role |
|----------|----------------|----------------|------|
| `docx-frontend` | 3000 | 3011 | Next.js frontend |
| `docx-backend` | 2000 | 2011 | Go (Fiber) backend |
| `docx-nginx` | 3011, 2011 | 3011, 2011 | Reverse proxy (frontend/backend entrypoints) |

---

### 5️⃣ Local URLs (Development/Test)

- **Frontend (UI):**  
  🔗 http://localhost:3011/docx-converter
  
- **Backend (API Base):**  
  🔗 http://localhost:2011/docx-converter-api
  
___

## 🧭 Directory Layout

```
docx-converter/
├── back/                  # Go backend
│   ├── cmd/server         # Main entry
│   ├── internal/          # Conversion logic
│   ├── go.mod / go.sum
│   └── .env
├── front/                 # Next.js frontend
│   ├── app/               # App router
│   ├── components/
│   ├── public/
│   └── .env
├── nginx/                 # NGINX configs
│   └── nginx.conf
└── docker-compose.yml
```

---

## 🧰 Troubleshooting

| Issue | Cause | Solution |
|-------|--------|----------|
| `404` after deploy | BasePath misconfigured | Check `NEXT_PUBLIC_APP_BASEPATH` and Nginx `location /docx-converter` |
| `CORS` error | Origin mismatch | Add correct origin in `CORS_ALLOW_ORIGINS` |
| `502 Bad Gateway` | Backend container not ready | Run `docker compose logs backend` |
| `pandoc` not found | Missing dependency in container | Ensure backend Dockerfile installs `pandoc` |
| Static assets 404 | BasePath mismatch | Confirm basePath in `next.config.ts` matches `/docx-converter` |

---

## 🏗️ Deployment Notes (For LE Admin)

Once deployed, the application will be served under:

- **Frontend (UI):**  
  🔗 https://www.lifelong.cmu.ac.th/docx-converter/

- **Backend (API Base):**  
  🔗 https://www.lifelong.cmu.ac.th/docx-converter-api/

All internal routing between the frontend, backend, and Nginx containers is handled automatically by Docker Compose.

---

## 🏫 Acknowledgement

This project was developed as a **Senior Project** for the  
**Department of Computer Engineering, Faculty of Engineering, Chiang Mai University.**

Developed by:  
**Jiradate Oratai**, **Nontapan Chanadee**, **Thatthana Sringoen**, and **Surapa Luangpiwdet**

Project Advisor:  
**Kampol Woradit**

in collaboration with the  
**Chiang Mai University School of Lifelong Education**,  
which serves as the primary stakeholder and future maintainer of this system.

© 2025 Chiang Mai University. All rights reserved.



