FROM golang:1.14-alpine as build

RUN apk add --no-cache git 

RUN wget https://github.com/tektoncd/cli/releases/download/v0.17.0/tkn_0.17.0_Linux_x86_64.tar.gz
RUN tar xvzf tkn_0.17.0_Linux_x86_64.tar.gz -C /usr/local/bin/ tkn

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

RUN apk add --no-cache git curl

RUN wget https://github.com/tektoncd/cli/releases/download/v0.17.0/tkn_0.17.0_Linux_x86_64.tar.gz
RUN tar xvzf tkn_0.17.0_Linux_x86_64.tar.gz -C /usr/local/bin/ tkn

RUN apk add util-linux fortune sed findutils openssh rsync \
      && curl -L https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl > /usr/bin/kubectl \
      && chmod +x /usr/bin/kubectl 

COPY --from=build /src/consumer /app/consumer

CMD [ "/app/consumer" ]