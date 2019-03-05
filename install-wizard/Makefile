DOCKER_IMAGE=kubermatic/install-wizard
CMD?=bash

install:
	npm install

start:
	./node_modules/.bin/ng serve

test-headless:
	./node_modules/.bin/ng test --watch=false --browsers=ChromeInDocker

build:
	./node_modules/.bin/ng build --prod --aot --build-optimizer --optimization

lint:
	./node_modules/.bin/tslint -c tslint.json -p src/tslint.json

# Helpers for devs not wanting to have a full webstack
# installed on their machines.

docker:
	docker build --build-arg "USER_ID=$$(id -u)" --build-arg "GROUP_ID=$$(id -g)" -t $(DOCKER_IMAGE) .

shell: docker
	docker run --rm -it -v $(realpath .):/home/node/app -w /home/node/app -p 4200:4200 $(DOCKER_IMAGE) $(CMD)
