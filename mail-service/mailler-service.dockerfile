
# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY maillerApp /app

CMD [ "/app/maillerApp" ]