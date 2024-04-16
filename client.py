import redis

# Connect to Redis on localhost (default port is 6379)
redis_client = redis.Redis(host='localhost', lib_name = None, lib_version = None)

# Try to get a sample key (likely doesn't exist)
sample_key = "my_sample_key"
value = redis_client.get(sample_key)

# Assert that the value is None (key doesn't exist)
assert value is None, f"Key '{sample_key}' unexpectedly has value: {value}"

# Set the key to 1 with an expiration of 5 seconds
key = "test_key"
redis_client.set(key, 1, ex=1)

# Get the value of the key we just set
value = redis_client.get(key)

# Assert that the value is 1
assert value == b'1', f"Key '{key}' has unexpected value: {value}"

# Wait for the key to expire (sleep for 6 seconds)
import time
time.sleep(2.1)

# Try to get the expired key
value = redis_client.get(key)

# Assert that the value is None (key expired)
assert value is None, f"Key '{key}' still has value after expiration: {value}"

print("All assertions passed! Script completed successfully.")

print(redis_client.execute_command("SET x"))

# GET x|@nil
# SET x|!ERR something wrong
# SET x 1|OK