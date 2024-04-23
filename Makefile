
swag:
	swag init -g ./internal/transport/http/server/server.go

deploy:
	./manual-deploy.sh

PHONY: swag, deploy

DEFAULT: swag
