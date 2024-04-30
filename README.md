# redtable

A redis server, backed by BigTable.

> CURRENT STATUS: PLAY THING

## Why

Redtable allows you to speak redis and persist in BigTable: it
acts as a proxy that translates redis / RESP commands into queries
to BigTable.

Since in many use-cases redis is purely used as a "simple" cache, but it's 
not as trivial to scale (single-threaded, with experiments to allow
multithreading, RAM size, etc etc), I've been eyeing moving some
workloads from caching on redis to BigTable, to ease operational overhead. 
As part of this fun thought-process, I've been wondering what a
proxy from RESP to BigTable would look like...

...and well, I didn't want to wonder for much longer.

The focus of redtable is to provide first-class support for commands
operating on a subset of data structures redis users might be familiar with:

* strings
* hashes
* sets
* sorted sets
* maps

## Supported commands

Currently, all these commands are supported:

```
APPEND
BITCOUNT
COPY
DBSIZE
DECR
DECRBY
DEL
ECHO
EXISTS
FLUSHALL
FLUSHDB
GET
GETDEL
KEYS
INCR
INCRBY
PING
RENAME
SET
SHUTDOWN
TIME
TTL
QUIT
```

These features are not supported (but most likely under evaluation):

```
BITCOUNT by BIT
GET does not return "WRONGTYPE Operation against a key holding the wrong kind of value" on the wrong type
KEYS with a pattern other than *
SET with EXAT
SET with PEXAT
SHUTDOWN with ABORT
```

Commands dealing with other data structures we're planning to support
are coming up...

## Local development

All you need is `docker`. Run `make` and see redtable come to life; you can then:

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
and run all of these commands one by one.

To run tests manually just:

```sh
# runs the services for local dev
make
# runs the tests
# first against an actual redis instance
# then same tests against redtable
make test
```

You can also test a single command with `make test cmd=$COMMAND_YOU_WANNA_TEST`:

```sh
$ make test cmd=del     
docker compose exec client python test.py redis 6379 del
# https://redis.io/docs/latest/commands/del/
@flushall
SET x 1|OK
SET x 1 > OK PASSED
SET y 1|OK
SET y 1 > OK PASSED
SET z 1|OK
SET z 1 > OK PASSED
DEL x y z a|$3
DEL x y z a > 3 PASSED

ALL TESTS PASSED
docker compose exec client python test.py redtable 6380 del
# https://redis.io/docs/latest/commands/del/
@flushall
SET x 1|OK
SET x 1 > OK PASSED
SET y 1|OK
SET y 1 > OK PASSED
SET z 1|OK
SET z 1 > OK PASSED
DEL x y z a|$3
DEL x y z a > 3 PASSED

ALL TESTS PASSED
```

## Considerations & misc

### Atomicity, transactions & what we do on a best-effort basis...

There are some predominant differences in how redis and
redtable operate under the hood as, for example, bigtable has sporadic
support for complex atomic/transactional operations. 

In general, it's really tough to compare
a distributed column / kv database such as BigTable to a (mostly)
single-threaded, in-memory DB like redis -- therefore things like
atomicity/transactions are implemented on a best-effort basis: if that is
problematic...then maybe redtable is not for you. 

A simple example is issuing a `DEL k1 k2 k3`:
redis returns the count of keys it deletes, but
bigtable does not support such operation (you can issue deletes in bulk,
but no way to know which rows existed), so we execute a bulk get
(to get the number of actual keys we're deleting) and then a bulk
delete, with obvious pitfalls -- here's some pseudo-code that
illustrates the process:

```sh
function del(keys):
    rows = bigtable.bulk_get(keys)
    bigtable.bulk_delete(keys)

    return len(rows)

# in process 0
set("x", 1)
set("y", 1)

# in process 1
del("x", "y", "z")

# in process 2
set("z", 1)
```

Depending on when process `1`/`2` are running, `1` may report
2 deletes while 3 have actually happened (`set` of `z` happens
between `1`'s bulk_get / bulk_delete).

Redtable will excel in those cases where you'd like to use redis 
as (mostly) a cache, but are struggling to scale -- as it combines
a very lightweight, multi-threaded redis-speaking server with the
low-latency of BigTable.

### Why not using BigTable's GC as a way to implement key expiration?

It could be possible to use [BigTable's GC](https://cloud.google.com/bigtable/docs/gc-cell-level)
to simulate redis' key expiry mechanism, but redtable instead stores the
expiry of a key into a separate bigtable column -- meaning rows could possibly
be fetched from bigtable, and "discarded" by the redtable server upon
figuring out that their expiry is in the past.

While the approach of manually checking the timestamp in a separate column
seems painful, the main reasons we do so is to be able to support bigtable's
atomic operations, namely append / increment -- these operations rely on
autmatically setting a server-generated timestamp when the cell is amended,
and using a `1s` GC policy would make it so you'd increment a counter, just to see
it disappear straight away.

### Data modeling

Currently, we store each kv pairs as BT rows with 1 exact column (`_values:value`),
but there's a case to be made for other examples of modeling, for example having one
row representing all strings, and have colums represent keys, with cells containing
their values.

There's still a lot to think about.

### Commands under evauation

These commands are currently not supported (but most likely under evaluation):

```
EXPIRE

BITFIELD
BITFIELD_RO
BITOP
BITPOS
BLMOVE
BLMPOP
BLPOP
BRPOP
BRPOPLPUSH
BZMPOP
BZPOPMAX
BZPOPMIN
DEBUG
DISCARD
DUMP
EXEC
EXPIREAT
EXPIRETIME
FAILOVER
FCALL
FCALL_RO
GETBIT
GETEX
GETRANGE
GETSET
HDEL
HELLO
HEXISTS
HGET
HGETALL
HINCRBY
HINCRBYFLOAT
HKEYS
HLEN
HMGET
HMSET
HRANDFIELD
HSCAN
HSET
HSETNX
HSTRLEN
HVALS
INCRBYFLOAT
INFO
LASTSAVE
LCS
LINDEX
LINSERT
LLEN
LMOVE
LMPOP
LOLWUT
LPOP
LPOS
LPUSH
LPUSHX
LRANGE
LREM
LSET
LTRIM
MGET
MIGRATE
MONITOR
MSET
MSETNX
MULTI
OBJECT
OBJECT|ENCODING
OBJECT|FREQ
OBJECT|HELP
OBJECT|IDLETIME
OBJECT|REFCOUNT
PERSIST
PEXPIRE
PEXPIREAT
PEXPIRETIME
PFADD
PFCOUNT
PFDEBUG
PFMERGE
PFSELFTEST
PSETEX
PSUBSCRIBE
PSYNC
PTTL
RANDOMKEY
READONLY
READWRITE
RENAMENX
REPLCONF
REPLICAOF
RESET
RESTORE
RESTORE-ASKING
ROLE
RPOP
RPOPLPUSH
RPUSH
RPUSHX
SADD
SAVE
SCAN
SCARD
SDIFF
SDIFFSTORE
SELECT
SETBIT
SETEX
SETNX
SETRANGE
SINTER
SINTERCARD
SINTERSTORE
SISMEMBER
SLAVEOF
SLOWLOG
SLOWLOG|GET
SLOWLOG|HELP
SLOWLOG|LEN
SLOWLOG|RESET
SMEMBERS
SMISMEMBER
SMOVE
SORT
SORT_RO
SPOP
SPUBLISH
SRANDMEMBER
SREM
SSCAN
SSUBSCRIBE
STRLEN
SUBSCRIBE
SUBSTR
SUNION
SUNIONSTORE
SUNSUBSCRIBE
SWAPDB
SYNC
TOUCH
TYPE
UNLINK
UNSUBSCRIBE
UNWATCH
WAIT
WAITAOF
WATCH
ZADD
ZCARD
ZCOUNT
ZDIFF
ZDIFFSTORE
ZINCRBY
ZINTER
ZINTERCARD
ZINTERSTORE
ZLEXCOUNT
ZMPOP
ZMSCORE
ZPOPMAX
ZPOPMIN
ZRANDMEMBER
ZRANGE
ZRANGEBYLEX
ZRANGEBYSCORE
ZRANGESTORE
ZRANK
ZREM
ZREMRANGEBYLEX
ZREMRANGEBYRANK
ZREMRANGEBYSCORE
ZREVRANGE
ZREVRANGEBYLEX
ZREVRANGEBYSCORE
ZREVRANK
ZSCAN
ZSCORE
ZUNION
ZUNIONSTORE
```

### Commands out of scope

```
ACL
ACL|CAT
ACL|DELUSER
ACL|DRYRUN
ACL|GENPASS
ACL|GETUSER
ACL|HELP
ACL|LIST
ACL|LOAD
ACL|LOG
ACL|SAVE
ACL|SETUSER
ACL|USERS
ACL|WHOAMI
ASKING
AUTH
BGREWRITEAOF
BGSAVE
CLIENT
CLIENT|CACHING
CLIENT|GETNAME
CLIENT|GETREDIR
CLIENT|HELP
CLIENT|ID
CLIENT|INFO
CLIENT|KILL
CLIENT|LIST
CLIENT|NO-EVICT
CLIENT|NO-TOUCH
CLIENT|PAUSE
CLIENT|REPLY
CLIENT|SETINFO
CLIENT|SETNAME
CLIENT|TRACKING
CLIENT|TRACKINGINFO
CLIENT|UNBLOCK
CLIENT|UNPAUSE
COMMAND
COMMAND|COUNT
COMMAND|DOCS
COMMAND|GETKEYS
COMMAND|GETKEYSANDFLAGS
COMMAND|HELP
COMMAND|INFO
COMMAND|LIST
CONFIG
CONFIG|GET
CONFIG|HELP
CONFIG|RESETSTAT
CONFIG|REWRITE
CONFIG|SET
CLUSTER
CLUSTER|ADDSLOTS
CLUSTER|ADDSLOTSRANGE
CLUSTER|BUMPEPOCH
CLUSTER|COUNT-FAILURE-REPORTS
CLUSTER|COUNTKEYSINSLOT
CLUSTER|DELSLOTS
CLUSTER|DELSLOTSRANGE
CLUSTER|FAILOVER
CLUSTER|FLUSHSLOTS
CLUSTER|FORGET
CLUSTER|GETKEYSINSLOT
CLUSTER|HELP
CLUSTER|INFO
CLUSTER|KEYSLOT
CLUSTER|LINKS
CLUSTER|MEET
CLUSTER|MYID
CLUSTER|MYSHARDID
CLUSTER|NODES
CLUSTER|REPLICAS
CLUSTER|REPLICATE
CLUSTER|RESET
CLUSTER|SAVECONFIG
CLUSTER|SET-CONFIG-EPOCH
CLUSTER|SETSLOT
CLUSTER|SHARDS
CLUSTER|SLAVES
CLUSTER|SLOTS
EVAL
EVALSHA
EVALSHA_RO
EVAL_RO
FUNCTION
FUNCTION|DELETE
FUNCTION|DUMP
FUNCTION|FLUSH
FUNCTION|HELP
FUNCTION|KILL
FUNCTION|LIST
FUNCTION|LOAD
FUNCTION|RESTORE
FUNCTION|STATS
GEOADD
GEODIST
GEOHASH
GEOPOS
GEORADIUS
GEORADIUSBYMEMBER
GEORADIUSBYMEMBER_RO
GEORADIUS_RO
GEOSEARCH
GEOSEARCHSTORE
LATENCY
LATENCY|DOCTOR
LATENCY|GRAPH
LATENCY|HELP
LATENCY|HISTOGRAM
LATENCY|HISTORY
LATENCY|LATEST
LATENCY|RESET
MODULE
MODULE|HELP
MODULE|LIST
MODULE|LOAD
MODULE|LOADEX
MODULE|UNLOAD
MOVE
MEMORY
MEMORY|DOCTOR
MEMORY|HELP
MEMORY|MALLOC-STATS
MEMORY|PURGE
MEMORY|STATS
MEMORY|USAGE
PUBLISH
PUBSUB
PUBSUB|CHANNELS
PUBSUB|HELP
PUBSUB|NUMPAT
PUBSUB|NUMSUB
PUBSUB|SHARDCHANNELS
PUBSUB|SHARDNUMSUB
PUNSUBSCRIBE
SCRIPT
SCRIPT|DEBUG
SCRIPT|EXISTS
SCRIPT|FLUSH
SCRIPT|HELP
SCRIPT|KILL
SCRIPT|LOAD
XACK
XADD
XAUTOCLAIM
XCLAIM
XDEL
XGROUP
XGROUP|CREATE
XGROUP|CREATECONSUMER
XGROUP|DELCONSUMER
XGROUP|DESTROY
XGROUP|HELP
XGROUP|SETID
XINFO
XINFO|CONSUMERS
XINFO|GROUPS
XINFO|HELP
XINFO|STREAM
XLEN
XPENDING
XRANGE
XREAD
XREADGROUP
XREVRANGE
XSETID
XTRIM
```

### No but for real, why?

![](https://imgflip.com/s/meme/Yao-Ming.jpg)

## Ack

Really, the heavylifting is done by [redcon](https://github.com/tidwall/redcon),
a beautiful project by [Josh Baker](https://tidwall.com/).