

gen_cert:
	rm -rf certs
	mkdir -p certs
	openssl genrsa -out ./certs/server.key 2048
	openssl req -new -x509 -key ./certs/server.key -out ./certs/server.pem -days 3650
