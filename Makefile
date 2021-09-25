all:
	echo "No default target defined."

deploy:
	make -C infrastructure init apply

destroy:
	make -C infrastructure destroy

app:
	make -C services/app build update

proxy:
	make -C services/proxy build update

things:
	make -C services/things build update

services: things

.PHONY: static deploy destroy app proxy services things
