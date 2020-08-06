FROM alpine:3.12

RUN apk add --no-cache ca-certificates

ADD ./k8s-jwt-to-vault-token /k8s-jwt-to-vault-token

ENTRYPOINT ["/k8s-jwt-to-vault-token"]
