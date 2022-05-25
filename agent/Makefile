IMG=tibbar/of-agent:v0.0.1
build:
	docker buildx build --platform="linux/amd64,linux/arm64" -t ${IMG} . --push
run:
	go run main.go start -f examples/example_flags.json
