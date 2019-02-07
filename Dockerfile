FROM golang:latest
ENV PKG_CONFIG_PATH=/usr/lib/pkgconfig
ENV GOPATH=/root/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$PATH:$GOBIN
ENV LD_LIBRARY_PATH=/usr/lib

RUN \
    git clone https://github.com/edenhill/librdkafka.git && \
    cd librdkafka && \
    ./configure --prefix /usr && \
    make && \
    make install

RUN \
    go get github.com/gorilla/mux && \
    go get -u github.com/confluentinc/confluent-kafka-go/kafka

ADD . $GOPATH/src/github.com/gkuhn1/event-driven-go
WORKDIR $GOPATH/src/github.com/gkuhn1/event-driven-go