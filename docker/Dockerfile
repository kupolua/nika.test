FROM alpine:latest

RUN apk update && \
    apk upgrade && \
    apk add --no-cache openssh git

COPY scripts/generate-site.bash / /

ENTRYPOINT ["/generate-site.bash"]
