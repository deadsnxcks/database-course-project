# syntax=docker/dockerfile:1

# ---------- build stage ----------
FROM node:22-alpine AS builder

WORKDIR /app

# Зависимости отдельным слоём: правка кода не инвалидирует npm ci
COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# tsc -b && vite build — сборка упадёт при ошибках типов, это нормально
RUN npm run build

# ---------- runtime stage ----------
FROM nginx:1.27-alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]