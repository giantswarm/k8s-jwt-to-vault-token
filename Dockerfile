FROM quay.io/giantswarm/docker-kubectl:1.18.2

RUN apk update \
  && apk upgrade \
  && apk add jq \
  && rm -rf /var/cache/apk/*

COPY k8s-jwt-to-vault-token /k8s-jwt-to-vault-token

ENTRYPOINT ["/k8s-jwt-to-vault-token"]
