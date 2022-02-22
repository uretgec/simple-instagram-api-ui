FROM alpine:latest

LABEL service="reposter"

RUN mkdir -p /usr/local/reweb

ARG ENV_FILE

COPY ${ENV_FILE} /usr/local/reweb/.env
COPY assets /usr/local/reweb/assets
COPY templates /usr/local/reweb/templates
COPY --from=reposter-builder /usr/local/go/bin/reweb /usr/local/reweb/reweb

WORKDIR /usr/local/reweb

ENTRYPOINT [ "./reweb",  "--env-file", ".env" ]