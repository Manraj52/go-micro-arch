# Notes
# https://security.stackexchange.com/questions/150078/missing-x509-extensions-with-an-openssl-generated-certificate
# https://stackoverflow.com/questions/68196502/failed-to-connect-to-a-server-with-golang-due-x509-certificate-relies-on-legacy

# Verify crt
# openssl x509 -in server.crt -noout -text

# Clean Up
echo "Cleaning up *.crt, *.csr, *.pem, *.key"
rm *.crt
rm *.csr
rm *.pem
rm *.key

# Summary
# Private files : ca.key, server.key, server.pem, server.crt
# Public files : ca.crt (for client), server.csr (for CA)

# Changes these CN's to match your hosts in your environment if needed.
SERVER_CN=go-micro-arch
EMAIL=test@test.com

# Step 1: Generate Certificate Authority(CA) + Trust Certificate (ca.crt)
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN=${SERVER_CN}"

# Step 2: Generate the Server Private Key (server.key)
openssl genrsa -passout pass:1111 -des3 -out server.key 4096

# Step 3: Get a certificate signing request from the CA (server.csr)
openssl req -new -subj "/CN=${SERVER_CN}" -key server.key -passin pass:1111 -out server.csr
# openssl req -new -subj "/CN=${SERVER_CN}" -key server.key -passin pass:1111 -addext "subjectAltName = DNS:localhost" -out server.csr
# openssl req -new -passin pass:1111 -subj "/CN=${SERVER_CN}" -addext "subjectAltName = DNS:localhost" -out server.csr

# Step 4: Sign the certificate with the CA created (self-signing) - server.crt
openssl x509 -req -extensions v3_ca -extfile ./ssl-extensions-x509.cnf -passin pass:1111 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt
# openssl ca -create_serial -cert server.crt -keyfile ca.key -days 365 -in server.csr -batch -out server.crt

# Step 5: Convert the server certificate to .pem format (server.pem) - for gRPC
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem
