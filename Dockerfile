FROM golang

ADD . $GOPATH/github.com/AkashGit21/hostelites

RUN go install github.com/AkashGit21/hostelites@latest

ENTRYPOINT $GOPATH/bin/hostelites

EXPOSE 8080
