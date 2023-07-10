FROM golang:1.21rc2-alpine3.18 as builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /yc-preemptible-killer

FROM alpine:3.14
COPY --from=builder /yc-preemptible-killer /usr/local/bin/yc-preemptible-killer
CMD ["yc-preemptible-killer"]