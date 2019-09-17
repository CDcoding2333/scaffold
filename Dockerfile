# Start from an golang:stretch image with the latest version of Go installed
FROM golang:stretch as build-env

LABEL maintainer="caodong2333@gmail.com"

WORKDIR /go/src/CDcoding2333/scaffold

# Copy the local package files to the container's workspace.
ADD . .

# Build the application using makefile
RUN CGO_ENABLED=0 \
	go build -a -tags netgo -installsuffix netgo -installsuffix cgo -ldflags '-w -s' \
	-o ./build/linux/scaffold CDcoding2333/scaffold

FROM alpine:3.7

COPY --from=build-env /go/src/CDcoding2333/scaffold/build/linux/scaffold /app/scaffold

# Run the app by default when the container starts
CMD /app/scaffold