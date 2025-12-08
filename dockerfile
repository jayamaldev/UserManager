FROM golang:latest AS builder

WORKDIR /user-manager

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /user-manager-app

EXPOSE 8080

CMD ["/user-manager-app"]