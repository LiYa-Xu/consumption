FROM golang:latest AS build
COPY . /go/build
WORKDIR /go/build
RUN go build -o consumption
FROM registry.access.redhat.com/ubi8/ubi-minimal
WORKDIR /root/
COPY --from=build /go/build/consumption .
EXPOSE 8080
CMD ["./consumption"]