FROM golang:1.17 as build
WORKDIR /opt/src
COPY . .
RUN groupadd -g 1000 appuser &&\
    useradd -m -u 1000 -g appuser appuser
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo /opt/src

FROM alpine:3.15
RUN apk add --no-cache git
LABEL "repository"="https://github.com/online-applications/version-action"
LABEL "version"="1.0.0"
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build --chown=1000:0 /opt/src/version-action /app

ENTRYPOINT [ "/app" ]