check_install:
	which swagger

swagger:check_install
	swagger generate spec -o ./swagger.yaml --scan-models
