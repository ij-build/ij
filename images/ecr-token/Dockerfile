FROM alpine:3.8

RUN apk -v --update add python py-pip && \
    pip install --upgrade pip && \
    pip install --upgrade awscli==1.14.5 && \
    apk -v --purge del py-pip && \
    rm /var/cache/apk/* && \
    mkdir ~/.aws

WORKDIR /ij
CMD ["/ij/get-ecr-token.sh"]
COPY get-ecr-token.sh /ij
