FROM golang:latest as gobuild

# Copy the local package files to the container's workspace.
#RUN mkdir /app

WORKDIR /go/src/project/
COPY . /go/src/project/

WORKDIR /go/src/project/cmd/samples/recipes/helloworld
RUN go build -o .


# This results in a single layer image
FROM scratch
COPY --from=gobuild /go/src/project/cmd/samples/recipes/helloworld /bin/project
COPY config/development.yaml /bin/project/config
RUN chmod a+rwx /bin/project
ENTRYPOINT ["/bin/project"]
CMD ["--help"]

EXPOSE 80 8080 8082