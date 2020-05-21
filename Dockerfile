# build image
FROM golang:1.12-alpine3.9 AS build-env

# install build tools
RUN apk update && apk upgrade && apk add --no-cache bash git openssh libc6-compat

# build
WORKDIR /app
# copy sources
COPY . .
# setup build env
WORKDIR /app/cli
# vendor build only can be executed outside the GOPATH
RUN go build -mod=vendor .

# distribution image
FROM alpine:3.9

# add CAs
# add libc6-compat
# https://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN apk --no-cache add ca-certificates libc6-compat

COPY --from=build-env /app/cli/cli .

# start
CMD ["./cli"]