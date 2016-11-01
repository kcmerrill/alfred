FROM golang:1.6
MAINTAINER kc merrill <kcmerrill@gmail.com>
RUN go get github.com/mitchellh/gox
COPY . /code
RUN go get github.com/kcmerrill/alfred
ENTRYPOINT "alfred"
