FROM golang:1.20 as builder

RUN mkdir build
COPY . /build

WORKDIR /build/cmd
RUN go mod download

RUN  go build -o main

CMD ["sh", "-c", "./main"]