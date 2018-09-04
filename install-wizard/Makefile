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
dockerbuild:
	docker build -t installer .
