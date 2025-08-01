FROM golang:1.22.5 AS go
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOBIN=/bin
RUN go install github.com/go-delve/delve/cmd/dlv@v1.8.2

FROM go AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/exclude-prefixes .

FROM build AS test
CMD go test -test.v ./...

FROM test AS debug
CMD dlv -l :40000 --headless=true --api-version=2 test -test.v ./...

FROM gcr.io/distroless/static-debian12:nonroot as runtime
COPY --from=build /bin/exclude-prefixes /bin/exclude-prefixes
USER 65532:65532
ENTRYPOINT ["/bin/exclude-prefixes"]