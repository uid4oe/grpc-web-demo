FROM golang:1.16.5-alpine as builder

EXPOSE 55052
RUN apk update && apk add git

WORKDIR /go/src/app
COPY . .
RUN go mod init
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app .

FROM scratch
WORKDIR /bin
COPY --from=builder /go/bin/app /app

ENTRYPOINT ["/app"]
