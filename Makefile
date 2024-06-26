run:
	docker compose up
test: test_redis test_redtable
test_redis:
	docker compose exec client python test.py redis 6379 ${cmd}
test_redtable:
	docker compose exec client python test.py redtable 6380 ${cmd}
unsupported:
	cat tests.txt | grep unsupported | awk '{split($$0,a,/[|]/); split(a[2],b,/(: )/); print b[2]}' | sort