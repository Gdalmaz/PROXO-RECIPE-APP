# App-Api-Cafe/Dockerfile

# syntax=docker/dockerfile:1.4
FROM golang:1.20.5-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOPATH=/go
ENV GOCACHE=/go-build

# Air'i kur
RUN go install github.com/cosmtrek/air@v1.44.0

# go.mod ve go.sum dosyalarını kopyala ve bağımlılıkları indir
COPY go.mod go.sum ./
RUN go mod download

# Tüm kaynak kodunu kopyala
COPY . .

# Çalıştırma komutu
CMD ["air", "-c", ".air.toml"]
