FROM ubuntu:22.04

USER root
RUN apt update -y
RUN apt install -y curl
RUN mkdir /app
WORKDIR /app

COPY . .

CMD ["/app/online_fashion_shop"]