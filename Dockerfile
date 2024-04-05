FROM golang:1.22.2-alpine

RUN apk update
RUN apk upgrade
RUN apk add nodejs
RUN apk add --no-cache bash
RUN node --version

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build ./cmd/app

RUN chmod +x ./app

EXPOSE 8080

CMD [ "./app" ]