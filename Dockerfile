FROM alpine:3.10.3

RUN apk add --update tzdata ca-certificates && \
    rm -rf /var/cache/apk/*

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

USER appuser
COPY dist/tides /app/tides
ENTRYPOINT ["/app/tides"]
