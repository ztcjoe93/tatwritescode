FROM golang:1.19.3-alpine

WORKDIR /twc-app

COPY ../../web/assets ./web/assets
COPY ../../web/templates ./web/templates

COPY cmd/twc-app/go.mod ./cmd/twc-app/go.mod
RUN cd cmd/twc-app && go mod download

COPY cmd/twc-app/*.go ./cmd/twc-app

RUN cd cmd/twc-app && go mod tidy
RUN cd cmd/twc-app && go build -o tatwritescode

EXPOSE 8080

CMD ./cmd/twc-app/tatwritescode