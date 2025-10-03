# ---------- Stage 1: Build ----------
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-indexer .


# ---------- Stage 2: Run ----------
FROM alpine:3.20

WORKDIR /app

# Create dedicated non-root user & group with fixed UID/GID
RUN addgroup -S appgroup && adduser -S appuser -G appgroup


# Copy binary & set strict ownership
COPY --chown=appuser:appgroup --from=builder /app/go-indexer /app/go-indexer

# Remove root entries so container only knows about appuser
RUN sed -i '/^root:/d' /etc/passwd \
    && sed -i '/^root:/d' /etc/group

# Drop privileges
USER appuser

# Run binary as entrypoint
ENTRYPOINT ["./go-indexer"]