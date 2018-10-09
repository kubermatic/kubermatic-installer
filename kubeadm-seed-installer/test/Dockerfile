FROM ruby:2-alpine3.7

WORKDIR /test/

ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

# naje abd g++ are there for native ruby extensions
# g++ should not be needed for pharos 1.3.x
RUN \
  apk update && \
  apk add go openssl unzip curl git musl-dev ca-certificates jq openssh-client bash gettext rsync make g++ && \
  curl -LO https://releases.hashicorp.com/vault/0.10.3/vault_0.10.3_linux_amd64.zip && \
  unzip vault_0.10.3_linux_amd64.zip && mv vault /usr/local/bin && \
  curl -LO https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_linux_amd64.zip && \
  unzip terraform_0.11.7_linux_amd64.zip && mv terraform /usr/local/bin && \
  mkdir -p /go/src /go/bin && \
  curl -L https://storage.googleapis.com/kubernetes-release/release/v1.10.5/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl && \
  chmod +x /usr/local/bin/kubectl && \
  mkdir -p /test && \
  go get -u -v github.com/heptio/sonobuoy

ENV PHAROS_VERSION=1.3.1
WORKDIR /tmp

RUN git clone https://github.com/kontena/pharos-cluster /tmp/pharos-cluster && \
    cd /tmp/pharos-cluster && \
    git reset --hard v${PHAROS_VERSION} && \
    gem build pharos-cluster.gemspec && \
    gem install pharos-cluster-${PHAROS_VERSION}.gem && \
    cd - && \
    rm -rf /tmp/pharos-cluster

COPY ./* /
COPY ./test /test

ENTRYPOINT ["/bin/bash"]
CMD [ "/test/run-test.sh" ]
