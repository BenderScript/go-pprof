FROM golang:1.15.3
RUN ls

SHELL ["/bin/bash", "-c"]
WORKDIR /go-pprof
COPY main.go .
COPY pprof.go .
COPY go.mod .
COPY go.sum .

RUN go build

CMD ./go-pprof
