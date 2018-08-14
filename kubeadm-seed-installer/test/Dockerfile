FROM alpine

WORKDIR /test/

ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN \
  apk update && \
  apk add go openssl unzip curl git musl-dev ca-certificates jq openssh-client bash gettext rsync && \
  curl -LO https://releases.hashicorp.com/vault/0.10.3/vault_0.10.3_linux_amd64.zip && \
  unzip vault_0.10.3_linux_amd64.zip && mv vault /usr/local/bin && \
  curl -LO https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_linux_amd64.zip && \
  unzip terraform_0.11.7_linux_amd64.zip && mv terraform /usr/local/bin && \
  mkdir -p /go/src /go/bin && \
  curl -L https://storage.googleapis.com/kubernetes-release/release/v1.10.5/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl && \
  chmod +x /usr/local/bin/kubectl && \
  mkdir -p /test && \
  go get -u -v github.com/heptio/sonobuoy

COPY ./* /
COPY ./test /test

ENTRYPOINT ["/bin/bash"]
CMD [ "/test/run-test.sh" ]