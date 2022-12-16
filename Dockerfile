FROM golang:1.19.3-alpine

WORKDIR /twc-app

COPY go.mod ./
RUN go mod download

COPY assets ./assets
COPY css ./css
COPY templates ./templates

COPY *.go ./

RUN go mod tidy
RUN go build -o tatwritescode

EXPOSE 8080

CMD ./tatwritescode