this cmd can verify whether certificate associate with private key

openssl pkcs12 -export -clcerts -in ssl.pem -inkey ssl.key -out ssl.p12
