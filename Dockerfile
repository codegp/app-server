FROM alpine
RUN mkdir /client
RUN mkdir /server

COPY client/dist /client/dist
COPY server/server /server/server
COPY server/templates /server/templates
WORKDIR /server

EXPOSE 8080

CMD ["./server"]
