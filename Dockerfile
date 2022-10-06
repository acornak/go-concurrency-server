# syntax=docker/dockerfile:1

# get go image
FROM golang:1.18-bullseye

# create working directory
WORKDIR /app

# copy files
COPY go.mod ./
COPY go.sum ./

# download dependencies
RUN go mod download

# copy scripts
COPY cmd ./cmd

# compile app
RUN go build -ldflags="-s -w" -o /server ./cmd/api

# expose port
EXPOSE 8080

#
CMD [ "/server" ]
