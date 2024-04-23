import redis, sys

script, host, port = sys.argv[:3]
command = None

if len(sys.argv) > 3:
    command = sys.argv[3]
    
    if command == "all":
        command = None

# Connect to Redis on localhost (default port is 6379)
redis_client = redis.Redis(host=host, port=port)
redis_client.flushall()

def is_command(line):
    needle = "# https://redis.io/docs/latest/commands/"
    ok = line.startswith(needle)
    cmd = None
    
    if ok:
        cmd = list(filter(None, line[len(needle):].split("/")))[0]
    
    return [ok, cmd]

current_command = None
with open("tests.txt") as file:
    for line in file:
        ok, cmd = is_command(line)
        
        if ok and command != None:
            current_command = cmd
            
        if current_command != command:
            continue
        
        line = line.rstrip("\n")
        print(line)
        
        if line in ("", " ") or line[0] == "#":
            continue
        
        if line.startswith("@sleep"):
            secs = line.split(" ")[1]
            import time
            time.sleep(float(secs))
            continue
        
        if line.startswith("@flushall"):
            redis_client.flushall()
            continue
        
        cmd, exp = line.split("|")
        
        # these ca be skipped, they're there as a reference of stuff to add
        if exp.startswith("!unsupported in redtable"):
            continue
        
        # int returns
        if exp[0] == "$":
            exp = int(exp[1:])
        # bool returns
        elif exp[0] == "^":
            exp = bool(int(exp[1:]))
        # null returns
        elif exp == "_":
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
            
        if type(got) is list:
            got.sort()
            got = str(got)
            
        assert exp == got, f'"{exp}" expected, got "{got}" (type {type(got)})'
        print(f"{cmd} > {exp} PASSED")
        
print("ALL TESTS PASSED")