# ---------- Stage 1: build ----------
FROM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

# ---------- Stage 2: run ----------
FROM alpine:3.22
WORKDIR /app
COPY --from=build /out/api .
EXPOSE 8080
CMD ["./api"]
