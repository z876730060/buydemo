FROM golang:1.26-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o buydemo .

FROM alpine:3.19

RUN apk add --no-cache sqlite-libs ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/buydemo .
COPY --from=builder /app/static ./static
RUN mkdir -p /app/data

EXPOSE 8080
VOLUME ["/app/data"]
CMD ["./buydemo"]
