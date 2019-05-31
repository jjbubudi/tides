ARG ARCH
FROM multiarch/alpine:${ARCH}-v3.9

RUN apk add --update tzdata && \
    rm -rf /var/cache/apk/*

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

USER appuser
COPY dist/tides /app/tides
ENTRYPOINT ["/app/tides"]
