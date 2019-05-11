ARG ARCH
FROM multiarch/alpine:${ARCH}-v3.9

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

USER appuser

COPY dist/tides /app/tides

CMD ["/app/tides"]
