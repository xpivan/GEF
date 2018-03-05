#!/bin/sh
########################################################
CERT_PATH=$HOME/.docker
CA_CERT=$CERT_PATH/ca.pem
CA_KEY=$CERT_PATH/ca-key.pem
CLIENT_CERT=$CERT_PATH/cert.pem
CLIENT_KEY=$CERT_PATH/key.pem
SERVER_CERT=$CERT_PATH/server.pem
SERVER_KEY=$CERT_PATH/server-key.pem
PASSPHRASE=pass

dhostname=egieudat
dip=$1

## Clean
sudo rm -f *.pem

## CA certificate generation
sudo openssl genrsa -aes256 -passout pass:$PASSPHRASE -out $CA_KEY 2048
sudo openssl req -new -x509 -days 365 -key $CA_KEY -sha256 -passin pass:$PASSPHRASE -subj "/C=FR/ST=MyState/O=MyOrg" -out $CA_CERT 

## Server certificate generation
sudo openssl genrsa -out $SERVER_KEY 2048 
sudo openssl req -subj "/CN=${dhostname}" -new -key $SERVER_KEY -out server.csr 2>/dev/null
echo subjectAltName = IP:${dip} > extfile.cnf
sudo openssl x509 -passin pass:$PASSPHRASE -req -days 365 -in server.csr -CA $CA_CERT -CAkey $CA_KEY -CAcreateserial -out $SERVER_CERT -extfile extfile.cnf

## Client certificate generation
sudo openssl genrsa -out $CLIENT_KEY 2048 
sudo openssl req -subj '/CN=client' -new -key $CLIENT_KEY -out client.csr 2>/dev/null
echo extendedKeyUsage = clientAuth > extfile.cnf
sudo openssl x509 -passin pass:$PASSPHRASE -req -days 365 -in client.csr -CA $CA_CERT -CAkey $CA_KEY -CAcreateserial -out $CLIENT_CERT -extfile extfile.cnf 

## Cleaning
sudo rm -f client.csr server.csr extfile.cnf ca.srl
sudo chmod 0400 $CA_KEY  
## !!!!!! WARNING !!!!!! The SERVER_KEY is supposed to have the 0400 mod but need to be set up on 0444 to allow scp transfer.
sudo chmod 0444 $CA_CERT $SERVER_CERT $CLIENT_CERT $SERVER_KEY $CLIENT_KEY
