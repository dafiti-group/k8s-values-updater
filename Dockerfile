# Build the manager binary
FROM golang as builder

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

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM alpine/git
# FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/k8s-values-updater .
RUN apk add --no-cache -u openssh
# USER nonroot:nonroot

# ENTRYPOINT ["/k8s-values-updater"]
# CMD ["bump"]
