.PHONY: test
test: /tmp/console/debug/data1/test1 /tmp/console/debug/data1/test2 /tmp/console/debug/data1/test3 /tmp/console/debug/data2/test1 /tmp/console/debug/data2/test2 /tmp/console/debug/data2/test3
	docker compose down
	docker compose up --build

/tmp/console/debug/%:
	mkdir -p '$(shell dirname '$@')'
	head -c 1000000 /dev/urandom > '$@'

