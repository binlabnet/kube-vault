FROM exelban/baseimage:golang-latest as build-app

WORKDIR /app/

COPY go.mod .
COPY go.sum .
COPY . .

ENV GO111MODULE=on

RUN go mod download
RUN go build -o ./bin/main


FROM exelban/baseimage:alpine-latest

MAINTAINER Serhiy Mitrovtsiy <mitrovtsiy@ukr.net>
COPY --from=build-app /app/bin /
ENTRYPOINT ./main