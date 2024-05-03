FROM postgres:latest

LABEL maintainer="Sato Kaito <satodeyannsu@gmail.com>"

RUN apt-get update \
    && apt-get install -y postgresql-contrib \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 5432
