FROM cgr.dev/chainguard/go:latest-dev AS builder

WORKDIR /wolfichef
ARG commitHash="null"
ENV GO111MODULE=on \
         CGO_ENABLED=0 \
         GOOS=linux
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -gccgoflags "-s -w" -o wolfichef
RUN apk update && apk add upx
RUN upx --lzma wolfichef
# You may choose to disable these two lines if you are not concerned about the binary size.

FROM cgr.dev/chainguard/busybox

WORKDIR /app
COPY --from=builder /wolfichef/wolfichef /app/wolfichef

EXPOSE 8000
ENTRYPOINT ["/app/wolfichef"]
