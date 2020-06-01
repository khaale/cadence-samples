FROM golang:latest as gobuild

# Copy the local package files to the container's workspace.
#RUN mkdir /app

WORKDIR /go/src/project/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY /cmd/samples/ /go/src/project/cmd/samples/

WORKDIR /go/src/project/cmd/samples/recipes/helloworld
RUN go build -o .


# This results in a single layer image
FROM debian
COPY --from=gobuild /go/src/project/cmd/samples/recipes/helloworld/helloworld /goapp
COPY /config/ /config/
RUN chmod -R 777 /goapp

USER appuser:appuser

#RUN ls -l /config
ENTRYPOINT ["/goapp"]
#CMD ["--help"]

EXPOSE 8080 8081