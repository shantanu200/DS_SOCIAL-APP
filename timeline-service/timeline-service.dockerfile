FROM alpine:latest

RUN mkdir /app

COPY timeLineApp /app

CMD [ "/app/timeLineApp" ]