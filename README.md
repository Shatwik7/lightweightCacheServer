# Redis-Compatible Cache Server (Reverse Engineered)

This project is a lightweight, reverse-engineered implementation of a Redis-compatible server. It supports basic caching operations such as `SET` and `GET`, and is designed to be compact and efficient. The server is compatible with the official Redis client, making it easy to integrate into existing systems that rely on Redis for caching.

## Features

- **Redis Protocol Compatibility**: The server implements a subset of the Redis protocol (RESP - REdis Serialization Protocol), allowing it to work seamlessly with the official Redis client.
- **Basic Caching Operations**:
  - `SET key value`: Stores a key-value pair in the cache.
  - `GET key`: Retrieves the value associated with a key.
- **LRU-Based Eviction**: Implemented a Least Recently Used (LRU) eviction policy to automatically remove old keys when memory limits are reached.
- **Lightweight and Efficient**: The server is designed to be minimalistic, focusing only on the essential caching functionality.
- **Concurrency Support**: The server uses a thread-safe key-value store, ensuring safe concurrent access to the cache.
- **Smaller Image**: Optimized the image size to 15mb

## How It Works

The server listens for incoming connections on a specified port (default: `:2345`). When a client connects, it processes commands sent in the Redis protocol format. The server supports the following commands:

1. **SET**: Stores a key-value pair in the cache.
   - Example: `SET mykey myvalue`
2. **GET**: Retrieves the value associated with a key.
   - Example: `GET mykey`
3. **HELLO**: A custom command that responds with server information.
   - Example: `HELLO world`
4. **CLIENT**: A placeholder command for future client-related functionality.

The server uses a simple in-memory key-value store (`KeyVal`) to manage the cache. It ensures thread safety using a read-write mutex (`sync.RWMutex`).

## Compatibility with Redis Client

This server is designed to be compatible with the official Redis client. You can use any Redis client library or the `redis-cli` tool to interact with this server. For example:

```bash
# Using redis-cli
$ redis-cli -p 2345
127.0.0.1:2345> SET mykey myvalue
OK
127.0.0.1:2345> GET mykey
"myvalue"
```



### Build and Run the Docker Container

```bash
docker build -t smallRedis .
```
```bash
docker run -p 2345:2345 --rm smallRedis
```
## Why This Optimization?

✅ Smaller Image: The final image contains only the compiled binary, reducing size.✅ Improved Security: No Go compiler or source code in production.✅ Faster Build & Deployment: Uses multi-stage builds to optimize caching.

