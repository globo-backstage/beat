# Settings

Backstage Beat is very configurable, all configurations are available below.

## Core

| Yaml path     | Enviroment variable | Default value  | Description |
| ------------- |---------------------| ---------------|-------------|
| log.level     | LOG_LEVEL           | info           | Application log level, are available: debug, info, warn, error, fatal and panic. |
| host          | HOST                | 0.0.0.0        | Host to serve API. |
| port          | PORT                | 3000           | Port to serve API. |
| database      | DATABASE            | mongo          | Database engine to use, are currently available: `mongo` and `redis` . |
| authentication| AUTHENTICATION      | static         | Authentication engine to use, is currently available: `static` . |


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
| redis.password         | REDIS_PASSWORD         |                | Password. |
