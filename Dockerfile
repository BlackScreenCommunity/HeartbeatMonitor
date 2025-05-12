FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git npm musl-dev gcc make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY package*.json ./
RUN npm install

COPY . .

RUN npm run bundle

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
ENV CC=gcc

ARG COMMIT_HASH
ARG VERSION_DATE_PART
ENV APP_VERSION=1.${VERSION_DATE_PART}-commit-${COMMIT_HASH}

RUN sed -i "s/#app-version#/${APP_VERSION}/g" ./internal/plugins/VersionPlugin.go

RUN go build -o HeartBeatMonitor .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/HeartBeatMonitor .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/internal/plugins/*/*.css ./templates/

ENTRYPOINT ["./HeartBeatMonitor"]
CMD ["-configFilePath", "/config/appsettings.json"]
