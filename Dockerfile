FROM golang:1.22
WORKDIR /mnt/homework
COPY . .
ENV CGO_ENABLED 0
RUN go build -o homework-object-storage ./cmd/server/...

# Docker is used as a base image so you can easily start playing around in the container using the Docker command line client.
FROM docker
COPY --from=0 /mnt/homework/homework-object-storage homework-object-storage
RUN apk add bash curl
CMD ["./homework-object-storage"]
