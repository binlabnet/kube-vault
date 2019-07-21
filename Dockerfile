FROM exelban/baseimage:golang-latest as build-app

WORKDIR /app/

COPY backend/go.mod .
COPY backend/go.sum .
COPY backend/. .

ENV GO111MODULE=on

RUN go mod download
RUN go build -o ./bin/main


FROM exelban/baseimage:node-latest as build-web

WORKDIR /app/

COPY frontend/package*.json ./
COPY frontend/yarn.lock ./
RUN npm install --silent

COPY frontend/ .
RUN yarn build:prod


FROM exelban/baseimage:alpine-latest

MAINTAINER Serhiy Mitrovtsiy <mitrovtsiy@ukr.net>
COPY --from=build-app /app/bin /
COPY --from=build-web /app/dist /dist
ENTRYPOINT ./main