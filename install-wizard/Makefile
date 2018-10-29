.PHONY: npminstall
npminstall: docker
	docker run --rm -v $(realpath .):/home/node/app -w /home/node/app install-wizard npm install

.PHONY: npmstart
npmstart: docker
	docker run --rm -it -v $(realpath .):/home/node/app -w /home/node/app --net host install-wizard npm start

.PHONY: shell
shell: docker
	docker run --rm -it -v $(realpath .):/home/node/app -w /home/node/app --net host install-wizard

.PHONY: prod
prod: docker
	docker run --rm -v $(realpath .):/home/node/app -w /home/node/app install-wizard ng build --prod --aot --build-optimizer --optimization

.PHONY: lint
lint: docker
	docker run --rm -v $(realpath .):/home/node/app -w /home/node/app install-wizard tslint -c tslint.json -p src/tslint.json

.PHONY: docker
docker:
	docker build --build-arg "USER_ID=$$(id -u)" --build-arg "GROUP_ID=$$(id -g)" -t install-wizard .
