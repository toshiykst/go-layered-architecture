FROM golang:1.19-alpine as builder

WORKDIR /go/src/github.com/toshiykst/go-layerd-architecture

COPY go.mod ./
#COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o go-layerd-architecture app/cmd/http/main.go


FROM gcr.io/distroless/static

COPY --from=builder /go/src/github.com/toshiykst/go-layerd-architecture/go-layerd-architecture app

EXPOSE 8080

CMD ["./app"]
