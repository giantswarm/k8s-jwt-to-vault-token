#!/bin/sh

vault_address=$1
role_name=$2
configmap=$3

token=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
echo "token: $token"
request=$(curl -k \
    --request POST \
    --data '{"jwt": "'"${token}"'", "role": "'"${role_name}"'"}' \
   ${vault_address}/v1/auth/kubernetes/login)

vault_token=$(echo ${request} | jq -r .auth.client_token)

kubectl create secret generic $configmap --from-literal=token=${vault_token}