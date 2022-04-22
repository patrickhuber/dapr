# this must be run from the root of the dapr source directory
# $ docker build -t gomemory:v1 -f pkg/sdk/examples/gomemory/Dockerfile .
FROM golang:1.18.1 as builder
WORKDIR /app/dapr/pkg/sdk/examples/gomemory
ADD  . /app/dapr
RUN go build -o dapr-gomemory-v1 .

# lts
FROM ubuntu:20.04
WORKDIR /app
EXPOSE 8888
COPY --from=builder /app/dapr/pkg/sdk/examples/gomemory/dapr-gomemory-v1 /app/
ENTRYPOINT [ "/app/dapr-gomemory-v1", "-p", "8888" ]