run:
	KV_VIPER_FILE=./config.local.yaml go run main.go run


migrate-up:
	KV_VIPER_FILE=./config.local.yaml go run main.go migrate up


migrate-down:
	KV_VIPER_FILE=./config.local.yaml go run main.go migrate down