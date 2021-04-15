FROM node:15

RUN apt-get update && \
    apt-get -y full-upgrade && \
    apt-get install -y curl && \
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"  && \
    unzip awscliv2.zip && \
    ./aws/install