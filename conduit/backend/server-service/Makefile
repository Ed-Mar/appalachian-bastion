.DEFAULT_GOAL := swagger

install_swagger:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger:
	@echo _________
	@echo Your gonna need Go Swagger for doc generate : go get -u github.com/go-swagger/go-swagger/cmd/swagger
	@echo or use this makefile with "$make install_swagger"

	swagger generate spec -o ./swagger.yaml --scan-models
	swagger generate spec -o ./swagger.json --scan-models