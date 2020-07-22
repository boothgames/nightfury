FROM golang:1.14 as build
ADD . /go/src/nightfury
WORKDIR /go/src/nightfury
RUN make compile-linux

FROM ubuntu
ENV METRICS_SERVER http://influxdb:8086
ENV METRICS_BUCKET nightfury
COPY --from=build /go/src/nightfury/out/nightfury /bin/nightfury
ENTRYPOINT nightfury server --metrics-server $METRICS_SERVER --metrics-bucket $METRICS_BUCKET