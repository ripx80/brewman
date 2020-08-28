FROM alpine:latest as builder

RUN apk update && apk upgrade && apk add --no-cache curl ca-certificates && \
    curl -L https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz --output upx-3.96-amd64_linux.tar.xz && \
    tar -xf upx-3.96-amd64_linux.tar.xz

FROM scratch
COPY --from=builder /upx-3.96-amd64_linux/upx /app
ENTRYPOINT [ "./app" ]

# docker run -it --rm -v $(pwd):/out ripx80/upx -9 -o out/brewman_upx out/brewman