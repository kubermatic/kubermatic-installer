FROM golang:1.13.8-alpine AS builder

RUN apk add -U git make

WORKDIR /home/kubermatic-installer/
COPY . .
RUN make && ./installer version

FROM alpine:3.11

ENV HELM_VERSION=v2.16.3
ENV KUBECTL_VERSION=v1.17.3

RUN apk add --no-cache ca-certificates

COPY --from=builder /home/kubermatic-installer/installer /usr/local/bin/

# add unprivileged user with a random UID
RUN adduser -h /home/kubermatic -u 8163 -D kubermatic
ENV HOME=/home/kubermatic

# run installer by default
ENTRYPOINT ["installer"]

# add Helm
RUN wget https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz && \
    tar xzf helm-${HELM_VERSION}-linux-amd64.tar.gz && \
    mv linux-amd64/helm /usr/local/bin && \
    rm -rf linux-amd64 helm-${HELM_VERSION}-linux-amd64.tar.gz && \
    helm version --short --client

# add kubectl
RUN wget -O /usr/local/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl && \
    chmod +x /usr/local/bin/kubectl && \
    kubectl version --short --client

# drop privileges
USER kubermatic
WORKDIR $HOME

# add remaining assets
COPY --chown=kubermatic:kubermatic charts ./charts
COPY --chown=kubermatic:kubermatic values.example.yaml .
