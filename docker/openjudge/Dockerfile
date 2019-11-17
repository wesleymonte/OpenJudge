FROM debian:stretch-slim

WORKDIR /service

EXPOSE 8080

COPY openjudge /service

RUN ["/bin/sh", "-c", "chmod +x openjudge"]

ENTRYPOINT ["./openjudge"]