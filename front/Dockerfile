# --- Build stage ---
FROM node:24-bookworm-slim AS build
WORKDIR /app

# รับ build args สำหรับค่าที่ต้องฝังใน bundle
ARG NEXT_PUBLIC_APP_BASEPATH
ARG NEXT_PUBLIC_BACKEND_URL
ENV NEXT_PUBLIC_APP_BASEPATH=$NEXT_PUBLIC_APP_BASEPATH
ENV NEXT_PUBLIC_BACKEND_URL=$NEXT_PUBLIC_BACKEND_URL

# ติดตั้ง deps
COPY package*.json ./
RUN npm ci

# สร้างแอป
COPY . .
RUN npm run build

# --- Runtime stage ---
FROM node:24-bookworm-slim AS run
WORKDIR /app
ENV NODE_ENV=production

# คัดลอกผล build แบบ standalone
COPY --from=build /app/.next/standalone ./
COPY --from=build /app/.next/static ./.next/static
COPY --from=build /app/public ./public

EXPOSE 3000
CMD ["node", "server.js"]
