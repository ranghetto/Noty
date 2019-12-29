FROM golang:1.13.5

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o noty

CMD ["./noty"]