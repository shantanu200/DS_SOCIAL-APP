FROM alpine:latest

RUN mkdir /app

COPY userRelationApp /app

CMD [ "/app/userRelationApp" ]
