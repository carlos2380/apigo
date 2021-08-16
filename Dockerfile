##
## Build
##

FROM golang:alpine AS build

WORKDIR /app

COPY . /app/

RUN go build -o apigo cmd/server/main.go
RUN go build -o client cmd/client/main.go

##
## Deploy
##

FROM alpine as apigo

WORKDIR /

COPY --from=build /app/apigo /apigo
CMD ["/apigo"]

FROM alpine as client

WORKDIR /

COPY --from=build /app/client /client
CMD ["/client"]