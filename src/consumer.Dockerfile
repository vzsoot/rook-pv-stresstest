FROM gcc:11.4.0 as builder

WORKDIR /workspace

COPY consumer.c consumer.c

RUN gcc -o consumer consumer.c && chmod a+x consumer

FROM rockylinux:9
WORKDIR /
COPY --from=builder /workspace/consumer .
USER 65532:65532

ENTRYPOINT ["/consumer"]
