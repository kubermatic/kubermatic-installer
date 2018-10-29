FROM ubuntu:18.04

ARG USER_ID=1000
ARG GROUP_ID=1000

RUN apt-get update -y
RUN apt-get install -y bash vim wget xz-utils make python git
RUN apt-get build-dep -y nodejs

# just get a recent nodejs with npm included
RUN wget -O /opt/n.tar.xz https://nodejs.org/dist/v10.4.1/node-v10.4.1-linux-x64.tar.xz
RUN tar x -C /opt -f /opt/n.tar.xz
RUN for x in /opt/node*/bin/*; do ln -sfn $x /usr/local/bin ; done

# polish environment
RUN groupadd -g ${GROUP_ID} node && \
    useradd -c 'Node User' -m -d /home/node -s /bin/bash -u ${USER_ID} -g ${GROUP_ID} node
USER node
ENV HOME /home/node
ENV PATH="${PATH}:/home/node/app/node_modules/.bin"
WORKDIR /home/node
