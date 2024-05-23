deploy:
	make -C infrastructure init apply

destroy:
	make -C infrastructure destroy

top:
	make -C services/top compile build clean

proxy:
	make -C services/proxy build

machine:
	make -C services/machine build

static:
	make -C services/static build

services: static machine

.PHONY: static deploy destroy top proxy services machine static
