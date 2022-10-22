FROM golang:1.16-alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ingestor examples/postgres-ingestion/ingestor/main.go
# RUN ls


FROM alpine:latest
WORKDIR /
COPY --from=build /app/ingestor /ingestor

ENV PORT=5000
EXPOSE 5000

ENTRYPOINT ["/ingestor"]
