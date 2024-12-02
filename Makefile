.PHONY: swagger-init

swagger-init:
	swag init --dir ./cmd,. --output ./docs