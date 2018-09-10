.PHONY: npminstall
npminstall:
	docker run --rm -v $(realpath .):/home/node/app -w /home/node/app installer npm install

.PHONY: npmstart
npmstart:
	docker run --rm -it -v $(realpath .):/home/node/app -w /home/node/app --net host installer npm start

.PHONY: shell
shell:
	docker run --rm -it -v $(realpath .):/home/node/app -w /home/node/app --net host installer

.PHONY: dockerbuild
prod:
	docker run --rm -v $(realpath .):/home/node/app -w /home/node/app installer ng build --prod --aot --build-optimizer --optimization

.PHONY: dockerbuild
dockerbuild:
	docker build --no-cache -t installer .
