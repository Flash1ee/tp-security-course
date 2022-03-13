#!/bin/sh
# Создание корневого сертификата - для сервера
mkdir -p ./certs/
openssl genrsa -out ./cert/ca.key 2048 # Создаем закрытый ключ
openssl req -new -x509 -days 3650 -key ./cert/ca.key -out ./cert/ca.crt -subj "/CN=Flash1ee CA" # создаем самоподписанный сертификат на основе закрытого ключа
openssl genrsa -out ./cert/cert.key 2048 # Создаем закрытый ключ
