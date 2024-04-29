FROM golang:1.22.1

WORKDIR /src

COPY . .

RUN go build -o main

RUN apt-get update && apt-get install -y curl

RUN echo 'HEALTHCHECK CMD curl -f http://rabbitmq:5672' >> Dockerfile

ENTRYPOINT [ "./main" ]
