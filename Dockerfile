FROM golang:1.24.2 AS go

ARG TARGETOS
ARG TARGETARCH

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOBIN=/bin
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go install github.com/go-delve/delve/cmd/dlv@v1.8.2

FROM go AS build
WORKDIR /build
COPY go.mod go.sum ./
COPY pkg ./pkg
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build ./pkg/imports
COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/exclude-prefixes .

FROM build AS test
CMD go test -test.v ./...

FROM test AS debug
CMD dlv -l :40000 --headless=true --api-version=2 test -test.v ./...

FROM alpine:3.21.3 AS runtime
COPY --from=build /bin/exclude-prefixes /bin/exclude-prefixes
ENTRYPOINT ["/bin/exclude-prefixes"]
