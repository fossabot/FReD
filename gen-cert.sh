#!/bin/bash

# usage: gen-cert.sh <name> <ip>
# check that we got the 2 parameters we needed or exit the script with a usage message
[ $# -ne 2 ] && { echo "Usage: $0 name ip"; exit 1; }

# give better names to parameter variables
NAME=$1
IP=$2

# generate a key
openssl genrsa -out "${NAME}".key 2048

# write the config file
# shellcheck disable=SC2086
cat > ${NAME}.conf <<EOF

[ req ]
default_bits = 2048
prompt = no
default_md = sha256
req_extensions = v3_req
distinguished_name = dn

[ dn ]
C = DE
ST = Berlin
L = Berlin
O = MCC
OU = FRED
EOF

# write the CN into the config file
echo "CN = ${NAME}" >> "${NAME}".conf

cat >> ${NAME}.conf <<EOF
[v3_req]
keyUsage = keyEncipherment, dataEncipherment, digitalSignature
extendedKeyUsage = serverAuth, clientAuth
subjectAltName = @alt_names

[alt_names]
IP.1 = 127.0.0.1
EOF

# write the IP SAN into the config file
echo "IP.2 = ${IP}" >> "${NAME}".conf

# generate the CSR
openssl req -new -key "${NAME}".key -out "${NAME}".csr -config "${NAME}".conf

# build the certificate
openssl x509 -req -in "${NAME}".csr -CA ca.crt -CAkey ca.key \
-CAcreateserial -out "${NAME}".crt -days 1825 \
-extfile "${NAME}".conf -extensions v3_req

# delete the config file and csr
rm "${NAME}".conf
rm "${NAME}".csr
