FROM golang:1.13.5

WORKDIR /app

ARG DATA_DIR=/app/data

RUN mkdir -p ${DATA_DIR}

COPY go.mod go.sum ./

RUN go mod download

COPY . .

VOLUME [${DATA_DIR}]

RUN go build -o noty

CMD ["./noty"]