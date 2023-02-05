FROM golang AS builder
WORKDIR /app

# Do a 'go mod download' before we copy over the Go code
# so we have a layer with our Go dependencies. This is a
# lazy way to have a cache so we don't beat up the sources.
# https://docs.docker.com/language/golang/build-images/
COPY go.mod go.sum ./
RUN go mod download -x

ADD . /app/
ENV GOTRACEBACK=all
ARG SKAFFOLD_GO_GCFLAGS
# RUN GOOS=linux GOARCH=amd64 go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o app cmd/platform/main.go
RUN GOOS=linux go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o app cmd/todo.go

# # FROM alpine
FROM debian:11 as server

COPY --from=builder /app/app /usr/local/bin/app
# RUN apk add --no-cache tzdata bash && \
#     addgroup -S dp_api && adduser -S -G dp_api -s /bin/bash dp_api
# USER dp_api
ENV GOTRACEBACK=all
ENTRYPOINT [ "/usr/local/bin/app", "server" ]
# ENTRYPOINT [ "/app/app" ]