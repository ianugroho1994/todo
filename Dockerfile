FROM alpine:latest

LABEL maintainer="ardiantonugroho <ardianto.nugroho1994@gmail.com>"

ARG http_port=1326

ENV TZ=Asia/Jakarta \
    PATH="/app:${PATH}"

RUN apk add --update --no-cache \
    sqlite \
    tzdata \
    libc6-compat \
    gcompat \
    ca-certificates \
    bash \
    && \
    cp --remove-destination /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone

# See http://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
# RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

ENV PORT $http_port
ENV IN_DB_HOST postgres
ENV IN_DB_PORT 5432
ENV IN_DB_NAME todo
ENV IN_DB_SSL disable
ENV IN_LISTEN_HOST todo-be
ENV IN_LISTEN_PORT $http_port
EXPOSE $http_port

WORKDIR /app

COPY ./build/todo /app/
RUN mkdir -p /app/cache
CMD ["./todo"]
