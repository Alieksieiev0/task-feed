FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/cosmtrek/air@latest

COPY . .
RUN rm -rf .git

CMD ["air", "-c", ".air.toml"]

