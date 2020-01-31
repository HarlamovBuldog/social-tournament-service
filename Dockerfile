FROM golang:alpine as builder

RUN apk update
RUN apk add ca-certificates

# Workdir is path in your docker image from where all your commands will be executed
WORKDIR /go/src/github.com/HarlamovBuldog/social-tournament-service

# Copy all from your project to WORKDIR
COPY . .

# Build the Go app.
# Rusulting files will be inside /go/bin/social-tournament-service
RUN make build

# Start a new build stage
FROM scratch

# Copy certificates from previous stage in order to make download from web work
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/social-tournament-service sts

# Starting bash  
ENTRYPOINT ["./sts"]
