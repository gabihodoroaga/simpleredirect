FROM alpine:3.13

WORKDIR app

ADD bin/simpleredirect .

RUN chmod +x simpleredirect

ENTRYPOINT ["/app/simpleredirect"]
