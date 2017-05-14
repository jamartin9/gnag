## Creating local test creds

# with golang
One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.

# With openssl
Generate a CA
1)    openssl req -out ca.pem -new -x509
        -generates CA file "ca.pem" and CA key "privkey.pem"

Generate server certificate/key pair
        - no password required.
2)    openssl genrsa -out server.key 1024
3)    openssl req -key server.key -new -out server.req
4)    openssl x509 -req -in server.req -CA CA.pem -CAkey privkey.pem -CAserial file.srl -out server.pem
        -contents of "file.srl" is a two digit number.  eg. "00"