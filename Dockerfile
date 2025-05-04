FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o AverageWatcherBot -ldflags="-s -w -extldflags '-static'"

FROM alpine

WORKDIR /app

COPY --from=builder /app/AverageWatcherBot .
COPY .env .
COPY ./storage ./storage

CMD [ "./AverageWatcherBot" ]