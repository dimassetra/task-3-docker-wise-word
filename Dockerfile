FROM golang:1.23.2-alpine

WORKDIR /app

RUN go version

COPY AUTHORS.md .
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o my-go-app

COPY LINKS.md .

EXPOSE 80

CMD ["./my-go-app"]