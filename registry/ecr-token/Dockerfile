FROM alpine:3.6

RUN apk -v --update add python py-pip && \
    pip install --upgrade pip && \
    pip install --upgrade awscli==1.14.5 && \
    apk -v --purge del py-pip && \
    rm /var/cache/apk/* && \
    mkdir ~/.aws

WORKDIR /ij
COPY get-ecr-token.sh /ij
CMD ["/ij/get-ecr-token.sh"]
