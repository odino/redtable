# redtable

A (somewhat) redis-compliant server, backed by BigTable.

## Why

(need to find a good excuse here)

## Local development

All you need is `docker`. Run `make` and see readtable come to life; you can then:

```sh
# on redis itself
$ nc localhost 6379
PING
+PONG
^C
# on redtable, magically behaving the same
$ nc localhost 6380
PING
+PONG
^C
```

## Tests

Our test suite is a simple series of commands and expected output
executed against a real redis instance, and then redtable, with the
expectations that both behave the same. We spawn a python redis client
and run all of this commands one by one.

To run tests manually just:

```sh
# runs the services for local dev
make
# runs the tests
# first against an actual redis instance
# then same tests against redtable
make test
```