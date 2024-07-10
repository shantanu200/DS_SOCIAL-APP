FROM alpine:latest

RUN mkdir /app

COPY notificationApp /app

CMD [ "/app/notificationApp"]
