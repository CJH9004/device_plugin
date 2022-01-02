.PHONY: build push

APP=virtual-device-plugin-example
IMG=cjh9004/$(APP):latest

build: 
	docker build -t $(IMG) .

push: build 
	docker push $(IMG)

run: build 
	docker run --rm -it $(IMG)

deployment: push 
	kubectl replace -f Deployment.yml --force
