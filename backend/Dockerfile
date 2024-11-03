
# syntax=docker/dockerfile:1

FROM golang:1.23.2-bookworm

WORKDIR /app

COPY . .

ENV GIN_MODE = 'release'
ENV PORT = '8080'

RUN go mod download
RUN go build -o /backend

EXPOSE 8080

CMD [ "/backend" ]
