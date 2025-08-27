FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o meserve ./cmd/meserve


FROM scratch


WORKDIR /app

COPY --from=builder /app/meserve .

EXPOSE 8080

ENTRYPOINT ["./meserve"]
