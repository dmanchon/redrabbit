FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

COPY . $GOPATH/src/github.com/dmanchon/redrabbit/
WORKDIR $GOPATH/src/github.com/dmanchon/redrabbit/

#get dependancies
#you can also use dep
RUN go get ./...

#build the binary
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -a -ldflags '-extldflags "-static"' -o /go/bin/app ./redrabbit/

# STEP 2 build a small image

# start from vanilla alpine
FROM alpine

# Copy our static executable
COPY --from=builder /go/bin/app /go/bin/app

# Default value for --host param
ENV REDIS="localhost:6379"

ENTRYPOINT ["/bin/sh", "-c", "/go/bin/app --host ${REDIS}"]
