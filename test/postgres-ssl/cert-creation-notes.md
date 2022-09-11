# Development Setup for SSL certs
Source: https://www.crunchydata.com/blog/ssl-certificate-authentication-postgresql-docker-containers

 ```sh
openssl req -new -x509 -days 10000 -nodes -out ca.crt \
-keyout ca.key \
-subj "/CN=root-ca" 

openssl req -new -nodes -out server.csr \
  -keyout server.key \
  -subj "/CN=localhost" \
  -addext "subjectAltName = DNS:localhost,IP:127.0.0.1"

openssl x509 -req -in server.csr -days 10000 \
    -CA ca.crt \
    -CAkey ca.key \
    -CAcreateserial \
    -out server.crt \
    -extfile <(printf "subjectAltName=DNS:localhost,IP:127.0.0.1")

rm server.csr

openssl req -new -nodes -out client.csr \
  -keyout client.key \
  -subj "/CN=mermerd_test"

openssl x509 -req -in client.csr -days 10000 \
      -CA ca.crt \
      -CAkey ca.key \
      -CAcreateserial \
      -out client.crt \
      -extfile <(printf "subjectAltName=DNS:localhost,IP:127.0.0.1")

rm client.csr

sudo chown 0:70 server.key
sudo chmod 640 server.key
```