
FROM alpine:latest

RUN mkdir /app

COPY tweetApp /app

CMD [ "/app/tweetApp"]
