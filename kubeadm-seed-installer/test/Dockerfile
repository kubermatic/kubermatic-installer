FROM alpine

RUN \
  apk update && \
  apk add unzip curl jq openssh-client bash gettext rsync && \
  curl -LO https://releases.hashicorp.com/vault/0.10.3/vault_0.10.3_linux_amd64.zip && \
  unzip vault_0.10.3_linux_amd64.zip && mv vault /usr/local/bin && \
  curl -LO https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_linux_amd64.zip && \
  unzip terraform_0.11.7_linux_amd64.zip && mv terraform /usr/local/bin
