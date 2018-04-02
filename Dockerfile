FROM golang:1.9
MAINTAINER kc merrill <kcmerrill@gmail.com>
COPY . /go/src/github.com/kcmerrill/alfred
WORKDIR /go/src/github.com/kcmerrill/alfred
RUN  go build -ldflags "-X main.Commit=`git rev-parse HEAD` -X main.Version=0.1.`git rev-list --count HEAD`" -o /usr/local/bin/alfred
ENTRYPOINT ["alfred"]
