FROM alpine
RUN mkdir /server

COPY server/server /server/server
COPY server/templates /server/templates
WORKDIR /server

EXPOSE 8080

CMD ["./server"]
