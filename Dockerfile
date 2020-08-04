FROM quay.io/giantswarm/docker-kubectl:1.18.2

RUN apk update \
  && apk upgrade \
  && apk add jq \
  && rm -rf /var/cache/apk/*

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
