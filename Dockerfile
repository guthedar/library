From golang:1.18.4

WORKDIR /home

COPY ./pkg /home

RUN cd /home && go build -o library

CMD ["/home/library"]
