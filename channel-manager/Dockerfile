FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-linkmode external -extldflags -static" -tags netgo -o /app/server

### 
FROM scratch 

COPY --from=build /app/server server

CMD ["/server"]
