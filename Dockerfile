# Stage 1: Build frontend
FROM node:22-alpine AS frontend
WORKDIR /app/ui
COPY ui/package.json ui/package-lock.json ./
RUN npm ci
COPY ui/ ./
RUN npm run build

# Stage 2: Build backend
FROM golang:1.25-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy the real frontend build into the embed directory
COPY --from=frontend /app/ui/dist ./internal/embed/dist
RUN CGO_ENABLED=0 go build -o /app/server ./cmd/server

# Stage 3: Production
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=backend /app/server .
EXPOSE 8080
ENTRYPOINT ["./server"]
