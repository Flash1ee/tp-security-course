mkdir certs/
openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt -subj "/CN=Flash1ee CA"
openssl genrsa -out cert.key 2048
sudo cp ca.crt /ust/local/share/ca-certificates/
sudo update-ca-certificates