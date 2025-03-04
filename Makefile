SWAG := swag
GO := go
SWAGGER_DIR := ./docs
SWAGGER_ARGS := init --dir ./,./cmd --generalInfo cmd/main.go --parseInternal

.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	$(SWAG) $(SWAGGER_ARGS)
	@echo "Swagger documentation generated in $(SWAGGER_DIR)"