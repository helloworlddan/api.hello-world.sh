all:
	app things proxy

deploy:
	make -C infrastructure init apply

destroy:
	make -C infrastructure destroy

app:
	make -C services/app compile build clean

proxy:
	make -C services/proxy build

things:
	make -C services/things build

services: things

.PHONY: static deploy destroy app proxy services things
