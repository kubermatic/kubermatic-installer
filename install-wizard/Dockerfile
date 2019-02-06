FROM ubuntu:18.04

ARG USER_ID=1000
ARG GROUP_ID=1000

RUN apt-get update -y
RUN apt-get install -y wget gnupg2

RUN echo 'deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main' | tee /etc/apt/sources.list.d/google-chrome.list
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add -

# re-enable sources since they got commented out somewhere in November 2018
RUN sed -i 's/# deb-src/deb-src/g' /etc/apt/sources.list

RUN apt-get update -y
RUN apt-get install -y bash vim wget xz-utils make python git google-chrome-stable
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
