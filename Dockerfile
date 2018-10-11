FROM scratch

COPY ./resizer-service /resizer-service

CMD ["/resizer-service"]
