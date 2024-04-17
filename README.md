# redtable

A (somewhat) redis-compliant server, backed by BigTable.

## Why

(need to find a good excuse here)

## Support

Supported commands:

```
APPEND
FLUSHALL
GET
SET
```

These features are not supported (but most likely under evaluation):

```
GET does not return "WRONGTYPE Operation against a key holding the wrong kind of value" on the wrong type
SET with EXAT
SET with PEXAT
```

You can generate this list with:

```sh
cat tests.txt | grep unsupported | awk '{split($0,a,/[|]/); split(a[2],b,/(: )/); print b[2]} | sort'
```

These commands are currently not supported (but most likely under evaluation):

```
ASKING
AUTH
BGREWRITEAOF
BGSAVE
BITCOUNT
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
COPY
DBSIZE
DEBUG
DECR
DECRBY
DEL
DISCARD
DUMP
ECHO
EVAL
EVALSHA
EVALSHA_RO
EVAL_RO
EXEC
EXISTS
EXPIRE
EXPIREAT
EXPIRETIME
FAILOVER
FCALL
FCALL_RO
FLUSHDB
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
GET
GETBIT
GETDEL
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
INCR
INCRBY
INCRBYFLOAT
INFO
KEYS
LASTSAVE
LATENCY
LATENCY|DOCTOR
LATENCY|GRAPH
LATENCY|HELP
LATENCY|HISTOGRAM
LATENCY|HISTORY
LATENCY|LATEST
LATENCY|RESET
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
MEMORY
MEMORY|DOCTOR
MEMORY|HELP
MEMORY|MALLOC-STATS
MEMORY|PURGE
MEMORY|STATS
MEMORY|USAGE
MGET
MIGRATE
MODULE
MODULE|HELP
MODULE|LIST
MODULE|LOAD
MODULE|LOADEX
MODULE|UNLOAD
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
SCRIPT
SCRIPT|DEBUG
SCRIPT|EXISTS
SCRIPT|FLUSH
SCRIPT|HELP
SCRIPT|KILL
SCRIPT|LOAD
SDIFF
SDIFFSTORE
SELECT
SETBIT
SETEX
SETNX
SETRANGE
SHUTDOWN
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
TIME
TOUCH
TTL
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