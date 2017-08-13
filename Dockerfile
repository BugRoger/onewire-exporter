FROM alpine:3.6
MAINTAINER Michael Schmidt <michael.j.schmidt@gmail.com>

ADD bin/linux/onewire-exporter onewire-exporter

ENTRYPOINT ["/onewire-exporter"]
