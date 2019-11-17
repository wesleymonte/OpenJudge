FROM debian:stretch-slim

WORKDIR /service

EXPOSE 8080

COPY open-judge /service

RUN ["/bin/sh", "-c", "chmod +x open-judge"]

ENTRYPOINT ["./open-judge"]