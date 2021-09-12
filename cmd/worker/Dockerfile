FROM golang:1.16-alpine

WORKDIR /app

# install curl
RUN apk --no-cache add curl

# include project libs
COPY pkg ./pkg

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# bring worker main
COPY cmd/worker/*.go ./

RUN go build -o /worker

CMD ["/worker"]
