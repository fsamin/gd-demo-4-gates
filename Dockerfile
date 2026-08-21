# A current Go, on the floating minor tag: the image gate scans this binary
# too, and a stale toolchain's own stdlib CVEs would fail it all by themselves.
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 go build -o /app .

# A real base image, so the running container has a shell to debug with.
FROM debian:11-slim
COPY --from=build /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
