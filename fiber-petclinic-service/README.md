# Fiber Pet Clinic Service

This example inherits every feature from gin-petclinic-service example. It has a circuit breaker middleware to handle
service failures gracefully.

## Features

- Middleware
    - Accepts JSON
    - Health Check
    - Monitor
    - Prometheus
    - Recovery
    - RequestID

## Endpoints

- `GET /api/pet-clinic/health`: Health check endpoint
- `GET /api/pet-clinic/metrics`: Prometheus metrics endpoint
- `GET /api/pet-clinic/monitor`: Monitor Page
- `GET /api/pet-clinic/info`: Application information endpoint
- `GET /api/pet-clinic/v1/owners/:id`: Get owner by ID

## Load Testing
- Using Grafana k6 for load testing.

## Fiber Zerolog


## [Circuit Breaker](https://docs.gofiber.io/contrib/circuitbreaker/)
### How It Works

1. Closed State:
    - Requests are allowed to pass normally.
    - Failures are counted.
    - If failures exceed a defined threshold, the circuit switches to Open state.

2. Open State:
    - Requests are blocked immediately to prevent overload.
    - The circuit stays open for a timeout period before moving to Half-Open.

3. Half-Open State:
    - Allows a limited number of requests to test service recovery.
    - If requests succeed, the circuit resets to Closed.
    - If requests fail, the circuit returns to Open.

### Circuit Breaker Benefits

* Prevents cascading failures in microservices.
* Improves system reliability by avoiding repeated failed requests.
* Reduces load on struggling services and allows recovery.

Install:

```bash
go get -u github.com/gofiber/fiber/v2
go get -u github.com/gofiber/contrib/circuitbreaker
```

Signature:

```go
circuitbreaker.New(config ...circuitbreaker.Config) *circuitbreaker.Middleware 
```

There are several configuration options available, such as:

1. Apply circuit breaker to all routes
2. Route and Route-Group specific circuit breaker
3. Circuit Breaker with custom failure handling
4. Circuit Breaker for External API calls or Outbound HTTP requests
5. Circuit Breaker with Concurrent Requests Handling
6. Circuit Breaker with Custom Metrics
7. Advanced: Multiple Circuit Breakers fo r Different Services

I recommend reading the [official documentation](https://docs.gofiber.io/contrib/circuitbreaker/) for more details and
examples. In this example, I have implemented a basic circuit breaker that applies to all routes.

## Live Reload
- Use [Air](https://github.com/air-verse/air) to enable live reloading during development. It automatically detects changes in your Go files and restarts the server, making development faster and more efficient.
- Use [Fiber Air](https://docs.gofiber.io/recipes/air/) example to configure the server with live reload



### Unit Tests
[Test Driven Design](UNIT_TEST.md) Documentation is available in the `UNIT_TEST.md` file. It includes instructions on
how to set up and run unit tests, as well as how to use Mockery for mocking dependencies.



