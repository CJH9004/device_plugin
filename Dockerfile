FROM golang:alpine as builder
WORKDIR /work
COPY go.mod go.mod
COPY main.go main.go
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN go build -o main .

FROM alpine:latest
COPY --from=builder /work/main /usr/bin/main
CMD [ "main" ]
