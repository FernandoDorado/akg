build: main.go
	arch -arm64 docker build --platform linux/amd64 -t bazhova/akg:latest .

push: main.go
	docker push bazhova/akg:latest
