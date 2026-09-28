.PHONY: test
test: /tmp/debug/data1/test1 /tmp/debug/data1/test2 /tmp/debug/data1/test3 /tmp/debug/data2/test1 /tmp/debug/data2/test2 /tmp/debug/data2/test3
	docker compose down
	docker compose up --build

/tmp/debug/%:
	mkdir -p '$(shell dirname '$@')'
	head -c 1000000 /dev/urandom > '$@'

