FROM golang:1.11.1-stretch

RUN mkdir static
RUN mkdir static/source
RUN mkdir static/resized

COPY ./static static

ADD resizer-service /

CMD ["/resizer-service"]
