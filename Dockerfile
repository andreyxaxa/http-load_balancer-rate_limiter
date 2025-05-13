# Step 1: Modules caching
FROM golang:1.24-alpine AS modules

WORKDIR /modules

COPY go.mod go.sum /modules/
RUN go mod download

# Step 2: Builder
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY --from=modules /go/pkg /go/pkg
COPY . /app

RUN CGO_ENABLED=0 go build -o /bin/app ./cmd/app

# Step 3: Final
FROM scratch

COPY --from=builder /app/config /config
COPY --from=builder /app/config.json /config.json
COPY --from=builder /bin/app /app

CMD ["/app", "/config.json"]