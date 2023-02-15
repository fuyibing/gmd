# GMD

Any application can publish messages as a producer, Any application can be used as
consumer consumption message. GMD maintains the subscription relationship of MQ, When
there is a message in the queue, GMD will deliver the message to the subscribers of
any application.

```shell
cd your_empty_path && \
git clone https://github.com/fuyibing/gmd.git . && \
go mod tidy && \
go build -o gmd && \
./gmd start
```

![Work flow](./docs/work-flow.png)

----

### Supports

1. [x] Aliyun MNS
2. [x] Rocket MQ
3. [ ] Rabbit MQ

----

### Guide

1. [HTTP API](./docs/api)
2. Utility
    1. Export Markdown documents. `go run main.go docs`
    2. Export Postman collection file. `go run main.go docs -a postman`
    3. Use docker container




