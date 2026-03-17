FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum main.go .
COPY src/ src/
RUN go run main.go
RUN ls -alR out

FROM scratch
COPY --from=builder --chown=65532:65532 /app/out/ .
