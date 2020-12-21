FROM golang:1.14-alpine as build

RUN apk add --no-cache git

WORKDIR /src 

RUN go get github.com/sirupsen/logrus
RUN go get github.com/streadway/amqp

RUN go get k8s.io/api/core/v1
RUN go get k8s.io/apimachinery/pkg/apis/meta/v1
RUN go get k8s.io/client-go/kubernetes
RUN go get k8s.io/client-go/rest

RUN go get k8s.io/api/apps/v1
RUN go get k8s.io/api/batch/v1




COPY consumer.go /src 

RUN go build consumer.go


FROM alpine as runtime

COPY --from=build /src/consumer /app/consumer

CMD [ "/app/consumer" ]