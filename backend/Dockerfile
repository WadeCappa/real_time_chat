
# syntax=docker/dockerfile:1

FROM golang:1.23.2-bookworm

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o /backend

# set to prod
RUN export GIN_MODE=release 

EXPOSE 8080

CMD [ "/backend" ]
