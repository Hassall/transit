FROM golang:1.16-alpine

WORKDIR /app

# include project libs
COPY pkg ./pkg

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# bring server main
COPY cmd/server/*.go ./

RUN go build -o /main

CMD ["/main"]
