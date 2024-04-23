FROM golang:1.22.1

WORKDIR /src

COPY . .

RUN go build -o main

ENTRYPOINT [ "./main" ]
