import redis, sys

script, host, port = sys.argv

# Connect to Redis on localhost (default port is 6379)
redis_client = redis.Redis(host=host, port=port)
redis_client.flushall()

with open("tests.txt") as file:
    for line in file:
        line = line.rstrip("\n")
        print(line)
        
        if line in ("", " ") or line[0] == "#":
            continue
        
        if line.startswith("@sleep"):
            print(line)
            secs = line.split(" ")[1]
            import time
            time.sleep(float(secs))
            continue
        
        cmd, exp = line.split("|")
        
        # these ca be skipped, they're there as a reference of stuff to add
        if exp.startswith("!unsupported in redtable"):
            continue
        
        # int returns
        if exp[0] == "$":
            exp = int(exp[1:])
            
        # null returns
        if exp == "_":
            exp = None
        
        try:
            got = redis_client.execute_command(cmd)
        except Exception as e:
            if exp[0] != "!":
                raise e
            
            exp = exp[1:]
            got = str(e)
        
        if type(got) is bytes:
            got = got.decode("utf-8")
            
        assert exp == got, f'"{exp}" expected, got "{got}" (type {type(got)})'
        print(f"{cmd} > {exp} PASSED")
        
print("ALL TESTS PASSED")