FROM gcc:11.4.0 as builder

WORKDIR /workspace

COPY producer.c producer.c

RUN gcc -o producer producer.c

FROM rockylinux:9
WORKDIR /
COPY --from=builder /workspace/producer .
USER 65532:65532

ENTRYPOINT ["/producer"]
