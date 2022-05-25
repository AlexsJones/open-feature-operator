openssl genrsa -out ca.key 2048  
openssl req -x509 -new -nodes -key ca.key -subj "/CN=192.168.1.116" -days 10000 -out ca.crt
openssl genrsa -out server.key 2048
cat << EOF >csr.conf
[ req ]
default_bits = 2048
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[ dn ]
CN = 192.168.1.116

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
IP.1 = 192.168.1.116

[ v3_ext ]
authorityKeyIdentifier=keyid,issuer:always
basicConstraints=CA:FALSE
keyUsage=keyEncipherment,dataEncipherment
extendedKeyUsage=serverAuth,clientAuth
subjectAltName=@alt_names
EOF
openssl req -new -key server.key -out server.csr -config csr.conf
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 10000 -extensions v3_ext -extfile csr.conf

mkdir -p $TMPDIR/k8s-webhook-server/serving-certs/
cp server.crt $TMPDIR/k8s-webhook-server/serving-certs/tls.crt
cp server.key $TMPDIR/k8s-webhook-server/serving-certs/tls.key

BUNDLE=`cat ca.crt | base64`
cat >> admissionwebhook.yaml << EOF
apiVersion: admissionregistration.k8s.io/v1beta1
kind: MutatingWebhookConfiguration
metadata:
  name: open-feature-mutating-admission-webhook
webhooks:
  - name: of-webhook
    clientConfig:
      service:
        name: open-feature-operator-webhook-service
        namespace: open-feature-operator-system
        path: "/mutate"
      caBundle: ${BUNDLE}
    rules:
      - operations: [ "CREATE","UPDATE" ]
        apiGroups: [""]
        apiVersions: ["v1"]
        resources: ["pods"]
EOF
kubectl create secret tls webhook-server-cert --key="server.key" --cert="server.crt" -n open-feature-operator-system             
