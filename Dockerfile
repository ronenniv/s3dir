FROM golang:1.17

RUN mkdir /s3dir
WORKDIR /s3dir

COPY . .

# Fetch dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download