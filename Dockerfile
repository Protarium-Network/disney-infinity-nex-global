# syntax=docker/dockerfile:1
ARG app_dir="/home/go/app"

FROM golang:1.27-alpine AS build
ARG app_dir
WORKDIR ${app_dir}

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o ${app_dir}/build/server .

FROM alpine:3.22 AS final
ARG app_dir
WORKDIR ${app_dir}

RUN addgroup go && adduser -D -G go go
USER go

COPY --from=build ${app_dir}/build/server ${app_dir}/server
CMD [ "./server" ]
