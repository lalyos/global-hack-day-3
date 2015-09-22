FROM golang:1.5.1

RUN go get -u github.com/jteeuwen/go-bindata/...
RUN go get github.com/progrium/go-shell

ENV TMP_DIR /kontainer
RUN mkdir -p $TMP_DIR/rootfs
#ADD https://github.com/sequenceiq/global-hack-day-3/releases/download/v0.3/runc-0.3.tgz $TMP_DIR/rootfs/
ADD ./runc-binary/release/runc $TMP_DIR/
ADD ./container.go /go/src/scratch/

ONBUILD ADD *.json $TMP_DIR/
ONBUILD ADD container.tar $TMP_DIR/rootfs/
ONBUILD RUN cd $TMP_DIR \
  && tar -cf kontainer.tar rootfs/ runc runtime.json config.json \
  && go-bindata kontainer.tar \
  && cp bindata.go /go/src/scratch/ \
  && cd /go/src/scratch/ \
  && go build -o kontainer .


