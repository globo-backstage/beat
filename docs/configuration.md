# Settings

Backstage Beat is very configurable, all configurations are available below.

## Core

| Yaml path     | Enviroment variable | Default value  | Description |
| ------------- |---------------------| ---------------|-------------|
| log.level     | LOG_LEVEL           | info           | Application log level, are available: debug, info, warn, error, fatal and panic. |
| host          | HOST                | 0.0.0.0        | Host to serve API. |
| port          | PORT                | 3000           | Port to serve API. |


## Database Engines

### MongoDB

| Yaml path      | Enviroment variable | Default value                        | Description |
| -------------- |---------------------| -------------------------------------|-------------|
| mongo.uri      | MONGO_URI           | localhost:27017/backstage_beat_local | Database URL. |
| mongo.user     | MONGO_USER          | None                                 | Username. |
| mongo.password | MONGO_PASSWORD      | None                                 | Password. |
| mongo.failFast | MONGO_FAILFAST      | true                                 | Cause connection and query attempts to fail faster when the server is unavailable. |


### Redis

| Yaml path              | Enviroment variable    | Default value  | Description |
| ---------------------- |----------------------- | ---------------|-------------|
| redis.host             | REDIS_HOST             | localhost:6379 | Host and port. |
| redis.db               | REDIS_DB               | 0              | Database. |
| redis.pool.maxIdle     | REDIS_POOL_MAXIDLE     | 10             | Maximum number of idle connections in the pool. |
| redis.pool.maxActive   | REDIS_POOL_MAXACTIVE   | 10             | Maximum number of connections allocated by the pool at a given time. When zero, there is no limit on the number of connections in the pool. |
| redis.pool.wait        | REDIS_POOL_WAIT        | true           | If true and the pool is at the `redis.pool.maxActive` limit, waits for a connection to be returned to the pool.|
| redis.pool.idleTimeout | REDIS_POOL_IDLETIMEOUT | 180e9          | Close connections after remaining idle for this duration. |
