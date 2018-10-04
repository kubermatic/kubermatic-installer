FROM alpine:3.8

ENV HELM_VERSION=v2.9.1
ENV KUBECTL_VERSION=v1.11.3
ENV KUBERMATIC_VALUES_YAML=/data/values.yaml
ENV KUBERMATIC_LISTEN_HOST=0.0.0.0

# add unprivileged user with a random UID
RUN adduser -h /home/kubermatic -u 8163 -D kubermatic
ENV HOME=/home/kubermatic

# run installer by default
ENTRYPOINT ["/home/kubermatic/installer"]

# add Helm
RUN wget https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz && \
    tar xzf helm-${HELM_VERSION}-linux-amd64.tar.gz && \
    mv linux-amd64/helm /usr/local/bin && \
    rm -rf linux-amd64 elm-${HELM_VERSION}-linux-amd64.tar.gz

# add kubectl
RUN wget -O /usr/local/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl && \
    chmod +x /usr/local/bin/kubectl

# drop privileges
USER kubermatic
WORKDIR $HOME

# add installer last
COPY installer .
COPY charts ./charts
COPY values.example.yaml .
