# TODO REST API Server

The server provides a REST API for managing tasks.

### Build:

```
git clone https://github.com/urvil38/todo-app-server.git
cd todo-app-server
make build
```

### Environment Variables:

| Name  |  Value  | Default Value  | Info
|:-----:|:---------------:|:--------------:|:---------|
|TODO_ADDRESS|0.0.0.0|0.0.0.0|Address on which server is running
|TODO_PORT|8080|8080|TCP port on which server is listening
|TODO_DEBUG_PORT|8081|8081|TCP port on which debug server is listening. Debug port should be different from the server port
|TODO_LOG_LEVEL|info, error, debug|info|LogLevel can be [info, debug, error, fatal]. If invalid log level is provided then info log level will be used as default
|TODO_LOG_FORMAT|text, json, json-pretty|text|LogFormat can be [json, json-pretty, text]
|JAEGER_AGENT_ENDPOINT|localhost:6831|""|TracingAgentURI instructs exporter to send spans to jaeger-agent at this address
|JAEGER_COLLECTOR_ENDPOINT|localhost:14268|""|TracingCollectorURI is the full url to the Jaeger HTTP Thrift collector

### Run Server:
```
./todo-app-server
```

## API

### Create Task:

```
curl --request POST \
  --url http://localhost:8080/task \
  --header 'Content-Type: application/json' \
  --data '{"task_name": "task1"}'
```

### Get Task:

```
curl --request GET \
  --url http://localhost:8080/task/1
```

### List Tasks:

```
curl --request GET \
  --url http://localhost:8080/tasks
```

### Update Task:

```
curl --request POST \
  --url http://localhost:8080/task/1 \
  --header 'Content-Type: application/json' \
  --data '{"task_name": "updated_task_1"}'
```

### Delete Task:

```
curl --request DELETE \
  --url http://localhost:8080/task/1
```

## Debug Server:

- The server records metrics and traces. These are available on [http://localhost:8081](http://localhost:8081)