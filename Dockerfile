FROM golang:1.8.1

ENV SRC_DIR $GOPATH/src/bitbucket.org/jungju/life
ENV BIN_FILE $GOPATH/bin/life
COPY . $SRC_DIR

WORKDIR $SRC_DIR
RUN go get github.com/tools/godep
RUN go build -o $GOPATH/bin/life *.go

ENV WEB_PORT 8373

CMD $BIN_FILE