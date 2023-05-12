FROM golang:1.18-alpine AS golang-image
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod .
RUN go mod download
#COPY *.go ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o b2c main.go
#============================================================
## deploy stage
FROM alpine as final
MAINTAINER "hugo_bh@yahoo.com"
LABEL service="img-b2c"
LABEL owner="hugo_bh"
RUN apk --no-cache add ca-certificates tzdata
#RUN apk --no-cache add redis
RUN apk add --no-cache redis
RUN mkdir /app
RUN chmod 777 /app
RUN mkdir /app/b2c
RUN chmod 777 /app/b2c
RUN mkdir /app/b2c/public
RUN chmod 777 /app/b2c/public
WORKDIR /app/b2c
COPY ./.env.local /app/b2c
COPY ./public/index.html /app/b2c/public
COPY ./public/listproduct.html /app/b2c/public
COPY ./public/listproductadmin.html /app/b2c/public
COPY ./public/login.html /app/b2c/public
COPY ./public/list.html /app/b2c/public
COPY ./public/customer.html /app/b2c/public
COPY ./public/carlist.html /app/b2c/public
COPY --from=golang-image /app/b2c /app/b2c
RUN chmod 777 /app/b2c/b2c
ENTRYPOINT ["/app/b2c/b2c"]
