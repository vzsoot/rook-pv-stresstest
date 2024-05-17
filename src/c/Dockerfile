FROM gcc:11.4.0 as builder

WORKDIR /workspace

COPY main.c main.c

RUN gcc -o application main.c && chmod a+x application

FROM rockylinux:9
WORKDIR /
COPY --from=builder /workspace/application .
USER 65532:65532

ENTRYPOINT ["/application"]
