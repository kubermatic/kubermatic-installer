FROM alpine:3.8

ENV HELM_VERSION=v2.9.1

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

# drop privileges
USER kubermatic
WORKDIR $HOME

# add installer last
COPY installer /home/kubermatic
