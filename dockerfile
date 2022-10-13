FROM golang:1.17 AS golang
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO11MODULE=on go build -a -o /main .

FROM alpine:latest
COPY --from=golang /main /uexporter
COPY --from=golang /app/plugins /plugins
RUN chmod +x /uexporter
