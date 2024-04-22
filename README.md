# redtable

A (somewhat) redis-compliant server, backed by BigTable.

## Why

(need to find a good excuse here)

The focus of redtable is to provide first-class support for commands
operating on a subset of data structures redis users might be familiar with:

* strings
* hashes
* sets
* sorted sets
* maps

## Supported commands


```
APPEND
BITCOUNT
COPY
DBSIZE
DEL
ECHO
FLUSHALL
FLUSHDB
GET
GETDEL
SET
SHUTDOWN
TTL
```

These features are not supported (but most likely under evaluation):

```
BITCOUNT by BIT
GET does not return "WRONGTYPE Operation against a key holding the wrong kind of value" on the wrong type
SET with EXAT
SET with PEXAT
SHUTDOWN with ABORT
```

You can generate this list with:

```sh
cat tests.txt | grep unsupported | awk '{split($0,a,/[|]/); split(a[2],b,/(: )/); print b[2]}' | sort
```

These commands are currently not supported (but most likely under evaluation):

```
DECR
DECRBY
EXISTS
EXPIRE
TIME
INCR
INCRBY
KEYS

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
COMMAND
COMMAND|COUNT
COMMAND|DOCS
COMMAND|GETKEYS
COMMAND|GETKEYSANDFLAGS
COMMAND|HELP
COMMAND|INFO
COMMAND|LIST
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
MOVE
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
PING
PSETEX
PSUBSCRIBE
PSYNC
PTTL
PUBLISH
PUBSUB
PUBSUB|CHANNELS
PUBSUB|HELP
PUBSUB|NUMPAT
PUBSUB|NUMSUB
PUBSUB|SHARDCHANNELS
PUBSUB|SHARDNUMSUB
PUNSUBSCRIBE
QUIT
RANDOMKEY
READONLY
READWRITE
RENAME
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

These commands are straight up not supported and out of scope for redtable, at least for now:

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
MEMORY
MEMORY|DOCTOR
MEMORY|HELP
MEMORY|MALLOC-STATS
MEMORY|PURGE
MEMORY|STATS
MEMORY|USAGE
SCRIPT
SCRIPT|DEBUG
SCRIPT|EXISTS
SCRIPT|FLUSH
SCRIPT|HELP
SCRIPT|KILL
SCRIPT|LOAD
```

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