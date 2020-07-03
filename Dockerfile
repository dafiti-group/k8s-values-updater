# Build the manager binary
FROM golang:alpine as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o k8s-values-updater .

#
FROM alpine
RUN apk add --no-cache -u openssh
WORKDIR /
COPY --from=builder /workspace/k8s-values-updater .

ENTRYPOINT ["/k8s-values-updater"]
CMD ["bump"]
