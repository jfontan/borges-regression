FROM golang:1.10-stretch

MAINTAINER jfontan

ENV LOG_LEVEL=debug

RUN apt-get update && \
    apt-get install -y dumb-init \
      git make bash gcc && \
    apt-get autoremove -y && \
    ln -s /usr/local/go/bin/go /usr/bin

ADD build/regression_linux_amd64/regression /bin/

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/bin/regression", "latest", "remote:master"]
