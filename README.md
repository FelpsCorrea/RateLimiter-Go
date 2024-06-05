# RateLimiter-Go

This project implements a rate limiter in Go, configurable to limit the maximum number of requests per second based on a specific IP address or an access token. The limiter logic is separated from the middleware, and the "limiter" information is stored and queried from a Redis database.

## Configuration
Create a .env file at the root of the project with the following variables:
```makefile
RATE_LIMIT_PER_IP=10
RATE_LIMIT_PER_TOKEN=100
BLOCK_DURATION=5m
REDIS_ADDR=redis:6379
```

To run the server, execute:
```sh
docker-compose up --build -d
```
To run the tests, execute:
```sh
docker-compose run --rm rate-limiter-tests
```

## Manual Testing
1. Start the server
    ```sh
    docker-compose up --build -d
    ```
2. Test the rate limit without using a token, considering we set a limit of 10 requests in the environment variables:
    ```sh
    for i in {1..10}; do curl -i http://localhost:8080/; done
    ```
    - The next request should be denied:
    ```sh
    curl -i http://localhost:8080/
    ```
3. Test the rate limit using a token, considering we set a limit of 100 requests in the environment variables:
    ```sh
    for i in {1..100}; do curl -i -H "API_KEY: abc123" http://localhost:8080/; done
    ```
    - The next request with the same token should be denied:
    ```sh
    curl -i -H "API_KEY: abc123" http://localhost:8080/
    ```
    - You can also test with a new token, which should be allowed:
    ```sh
    curl -i -H "API_KEY: abc321" http://localhost:8080/
    ```