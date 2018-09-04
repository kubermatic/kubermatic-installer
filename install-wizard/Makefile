.PHONY: npminstall
npminstall: dockerbuild
	docker run --rm -v $(realpath .):/home/node/app -w /home/node/app installer npm install

.PHONY: npmstart
npmstart: dockerbuild
	docker run --rm -it -v $(realpath .):/home/node/app -w /home/node/app --net host installer npm start

.PHONY: npmstart
shell: dockerbuild
	docker run --rm -it -v $(realpath .):/home/node/app -w /home/node/app --net host installer

.PHONY: dockerbuild
dockerbuild:
	docker build -t installer .
