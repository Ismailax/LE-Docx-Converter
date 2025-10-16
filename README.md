# CMU Lifelong Education — Course Document Converter

A web application for extracting and displaying structured course information from `.docx` files.  
Developed for **Chiang Mai University School of Lifelong Education**.

---

## 🧩 Project Structure

| Directory | Description |
|------------|-------------|
| `back/` | Go (Fiber) backend API for document conversion |
| `front/` | Next.js frontend for file upload and parsed data display |

---

## 🚀 Getting Started (Development)

### 1. Clone the repository

```bash
git clone https://github.com/Ismailax/docx-converter-demo.git
cd docx-converter-demo
```

---

### 2. Backend Setup (`back/`)

#### 📋 Prerequisites
- Go 1.21 or newer  
- Docker (required for Pandoc conversion)

#### ⚙️ Environment Variables
Create a `.env` file inside the `back/` directory using the provided example:

```bash
cp back/.env.example back/.env
```

Then edit the file as needed.  
Example variables:

```bash
PORT=
CORS_ALLOW_ORIGINS=
MAX_UPLOAD_MB=
```

> In production, make sure to set these values according to your deployment environment.

#### 🛠️ Installation & Run

Pull the Pandoc Docker image (used for conversion):

```bash
docker pull pandoc/core:latest
```

Then start the backend:

```bash
cd back
go run ./cmd/server
```

The backend server will start at the port defined in your `.env` file.

---

### 3. Frontend Setup (`front/`)

#### 📋 Prerequisites
- Node.js v18+ or newer  
- npm, yarn, or pnpm

#### ⚙️ Environment Variables
Create a `.env.local` file inside the `front/` directory:

```bash
cp front/.env.example front/.env.local
```

Then edit the values to match your environment.

Example variables:

```bash
NEXT_PUBLIC_APP_BASEPATH=
NEXT_PUBLIC_BACKEND_URL=
```

#### 📦 Installation

```bash
cd front
npm install
```

#### 🧭 Development Server

```bash
npm run dev
```

The frontend will start at the port defined in your script (default: 3000).

---

## 🖥️ Usage

1. Open the frontend in your browser.
2. Upload a `.docx` course document.
3. The extracted course information will be displayed.

---

## ⚠️ Troubleshooting

| Issue | Solution |
|--------|-----------|
| **Pandoc errors** | Ensure Docker Desktop is running and the `pandoc/core:latest` image is available. |
| **CORS errors** | Check that `CORS_ALLOW_ORIGINS` in backend `.env` matches your frontend URL. |
| **Port conflicts** | Change port values in the `.env` files if needed. |
| **TinyMCE assets missing** | Ensure that `public/tinymce/` exists in the frontend build output. |

---

## 🏗️ Deployment Notes

| Component | Container Port | Public URL (via reverse proxy) |
|------------|----------------|--------------------------------|
| Frontend | `3000` | `/docx-converter/` |
| Backend (API) | `2000` | `/docx-converter-api/` |

Both components are configured entirely via `.env` files — **no code modification required**.

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
