.PHONY: build push

APP=virtual-device-plugin
IMG=cjh9004/$(APP):latest

build: 
	docker build -t $(IMG) .

push: build 
	docker push $(IMG)

run: build 
	docker run --rm -it $(IMG)

daemonset: push 
	kubectl replace -f DaemonSet.yml --force
