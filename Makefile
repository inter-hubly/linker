# Makefile para Build Docker

# Definir caminhos de chave SSH
SSH_KEY_PATH=/home/saimon/.ssh/id_ed25519_no_passphrase

# Build padrão (docker)
build-docker:
	@SSH_KEY=$(shell cat $(SSH_KEY_PATH) | base64 -w 0) && \
	docker build --build-arg SSH_KEY="$${SSH_KEY}" --build-arg ENVIRONMENT=docker -t hubly:docker . && \
	docker tag hubly:docker saimonribeiros/hubly:docker && \
	docker push saimonribeiros/hubly:docker

# Build para desenvolvimento (dev)
build-dev:
	@SSH_KEY=$(shell cat $(SSH_KEY_PATH) | base64 -w 0) && \
	docker build --build-arg SSH_KEY="$${SSH_KEY}" --build-arg ENVIRONMENT=dev -t hubly:dev . && \
	docker tag hubly:dev saimonribeiros/hubly:dev && \
	docker push saimonribeiros/hubly:dev
