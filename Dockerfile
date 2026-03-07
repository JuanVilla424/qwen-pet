FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o pet ./cmd/pet

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /build/pet /pet
COPY --from=builder /build/configs /configs

VOLUME /data

ENTRYPOINT ["/pet"]
