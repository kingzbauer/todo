FROM golang:1.13.5-stretch as builder

MAINTAINER Jack Kinga <jackmwangi@gmail.com>

WORKDIR /home/app/

COPY ./go.* ./

RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app

FROM scratch

MAINTAINER Jack Kinga <jackmwangi@gmail.com>

WORKDIR /usr/bin
COPY --from=builder /home/app/app ./app

EXPOSE 8000

ENTRYPOINT ["app"]
