# redtable

A (somewhat) redis-compliant server, backed by BigTable.

## Why

(need to find a good excuse here)

## Tests

```sh
# runs the services for local dev
make
# runs the tests, first against an actual redis instance, then same tests against redtable
make test
```