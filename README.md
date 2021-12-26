# gateway-queue

a centralized Discord gateway identify ratelimiter. similar to twilight's [gateway-queue](https://github.com/twilight-rs/gateway-queue).

set `CONCURRENCY` env variable if you have a raised `max_concurrency` (if you're on large bot sharding).

gateway-queue binds to 0.0.0.0:8080 by default; set the `ADDR` env variable (with format `[host]:<port>`) to change that.
