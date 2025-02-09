# pull the base image
FROM golang:latest AS builder
# create base working directory inside container

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-project .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/go-project .
ENTRYPOINT ["./go-project"]
CMD ["appscode-api-server"]