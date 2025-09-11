# CMU Lifelong Education Course Document Converter Demo

A web application for extracting and displaying course information from `.docx` files.

## Project Structure

- `back/` — Go (Fiber) backend API for document conversion
- `front/` — Next.js frontend for file upload and data display

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/Ismailax/docx-conveter-demo.git
cd docx-conveter-demo
```
### 2. Backend Setup (back/)

#### Prerequisites

- Go 1.21 or newer
-	Docker (required for Pandoc conversion)

#### Environment Variables

Create a .env file in the back/ directory.

```bash
PORT=8080
FRONTEND_URL=http://localhost:3000
MAX_UPLOAD_MB=10
```

#### Installation & Run

Before running the backend, make sure to pull the Pandoc Docker image:

```bash
docker pull pandoc/core:latest
```

Then start the backend:

```bash
cd back
```

**Windows**
```bash
go run .\cmd\server
```

**macOS / Linux**
```bash
go run ./cmd/server
```

The backend server will start at http://localhost:8080.

> **Note:** Ensure that Docker Desktop is running before starting the backend, as Pandoc is executed inside a temporary Docker container created from the pandoc/core:latest image.

### 3. Frontend Setup (front/)

#### Prerequisites

- Node.js v18+ (or compatible)
- npm, yarn, or pnpm
 
#### Environment Variables

Create a .env file in the front/ directory.

```bash
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
```

#### Installation

```bash
cd front
npm install          # or: yarn install, pnpm install
```

#### Development Server

```bash
npm run dev          # or: yarn dev, pnpm dev
```

The frontend server will start at http://localhost:3000.

## Usage

1.	Open the frontend in your browser.
2.	Upload a .docx course document.
3.	The extracted course information will be displayed on the page.

## Troubleshooting

- **Pandoc errors:** Ensure Docker Desktop is running.
- **CORS errors:** Check that the backend .env value for FRONTEND_URL matches your frontend URL.
- **Port conflicts:** If default ports are in use, update the port settings in the backend and frontend.




