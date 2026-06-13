# Lunar Deploy Agent Testing Checklist

---

## CLI Tests

### Validate config

```bash
go run ./cmd/lunar-agent validate --config config.example.yaml
```
Expected:
* Config loads successfully
* Config validation passes

### List deployments
```bash
go run ./cmd/lunar-agent list --config config.example.yaml
```
Expected:
* Configured deployments are listed
* Repository path and step count are shown

### Run deployment
```bash
go run ./cmd/lunar-agent deploy --deployment happy-path --events --event-json --json
```
Expected:
* Deployment runs
* Events are printed
* JSON result is printed
* Exit code is ``0``

### Run failing deployment
```bash
go run ./cmd/lunar-agent deploy --deployment failure-test --json
echo $?
```
Expected
* Deployment fails
* JSON result is printed
* Exit code is ``1``

### View History
```bash
go run ./cmd/lunar-agent history
go run ./cmd/lunar-agent history --json
go run ./cmd/lunar-agent history --summary
```
Expected:
* History entries are readable
* JSON output works
* Summary output works

### View last deployment
```bash
go run ./cmd/lunar-agent last
go run ./cmd/lunar-agent last --json
```
Expected:
* Most recent deployment is displayed

### View local status
```bash
go run ./cmd/lunar-agent status
go run ./cmd/lunar-agent status --json
```
Expected:
* Agent status is displayed

---

## HTTP API Tests

Start server:
```bash
go run ./cmd/lunar-agent serve
```

Set token:
```bash
TOKEN="local-dev-token"
```

### Public health endpoint
```bash
curl http://localhost:8080/health
```
Expected:
* Health report returns
* Status is healthy or unhealthy with check details

### Protected status endpoint
```bash
curl http://localhost:8080/status \
  -H "Authorization: Bearer $TOKEN"
```
Expected:
* Status JSON returns

### List deployments
```bash
curl http://localhost:8080/deployments \
  -H "Authorization: Bearer $TOKEN"
```
Expected:
* Deployment list returns

### Deployment History
```bash
curl http://localhost:8080/history \
  -H "Authorization: Bearer $TOKEN"
```
Expected:
* Deployment history returns

### Queue status
```bash
curl http://localhost:8080/queue \
  -H "Authorization: Bearer $TOKEN"
```
Expected:
* Queue summary returns

### Node metadata
```bash
curl http://localhost:8080/queue \
  -H "Authorization: Bearer $TOKEN"
```
Expected:
* Node metadata returns

### Heartbeat
```bash
curl http://localhost:8080/heartbeat \
  -H "Authorization: Bearer $TOKEN"
```
Expected:
* Heartbeat response returns uptime, queue summary, and node identity

### Trigger deployment
```bash
curl -X POST http://localhost:8080/deploy \
  -H "Content-Type: applicatioin/json" \
  -H "Authorization: Bearer $TOKEN"
  -d '{"deployment":"happy-path"}'
```
Expected:
* Response status is ``202 Accepted``
* Job object is returned

### List jobs
```bash
curl http://localhost:8080/jobs \
  -H "Authorization: Bearer $TOKEN"
```
Expected:
* Job list returns

### Get Job by ID
```bash
curl http://localhost:8080/jobs/<job-id> \
  -H "Authorization: Bearer $TOKEN"
```
Expected:
* Specific job returns

### Cancel job
```bash
curl -X DELETE http://localhost:8080/jobs<job-id> \
  -H "Authorization: Bearer $TOKEN"
```
Expected:
* Job is marked canceled if cancellable

---

## Swagger UI Tests

Start server:
```bash
go run ./cmd/lunar-agent serve
```

Open:
```ignorelang
http://localhost:8080/swagger/index.html
```
Expected:
* Swagger UI loads
* OpenAPI spec loads
* API endpoints are visible

---

## WebSocket Tests

Open browswer console from:

```ignorelang
http://localhost:8080/swagger/index.html
```
Run:
```javascript
const socket = new WebSocket("ws://localhost:8080/events");

socket.onopen = () => console.log("connected");

socket.onmessage = (event) => {
    console.log("DEPLOY EVENT:", JSON.parse(event.data));
}

socket.onerror = (error) => console.error("socket error", error);

socket.onclose = (event) => console.log("closed", event.code, event.reason);
```

Trigger deployment:
```bash
curl -X POST http://localhost:8080/deploy \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN"
  -d '{"deployment":"happy-path"}'
```
Expected:
* Browser receives job, deployment, and step events

---

## Real PM2 Deployment Test
Use a real node test app with:
```json
{
  "scripts": {
    "start": "node server.js",
    "build": "echo Build completed"
  }
}
```

Deploy using:
```yaml
real-node-pm2-test:
  repository_path: "/absolute/path/to/test-node-app"
  preset: "node_pm2"
  processName: "lunar-test-app"
  health_check_url: "http://localhost:3000"
  health_check_expected_status: 200
```

Run:
```bash
go run ./cmd/lunar/agent deploy --deployment real-node-pm2-test --events --event-json --json
```
Expected:
* Git pull runs
* npm install runs
* PM2 starts or restarts app
* PM2 saves process list
* PM2 status passes
* HTTP health check passes

Verify:
```bash
pm2 list
curl http://localhost:3000
```

---

## Failure Tests

### Missing deployment
```bash
go run ./cmd/lunar-agent deploy --deployment does-not-exist
```
Expected:
* Deployment fails clearly

### Bad config
Temporarily break ``config.example.yaml``, then run:
```bash
go  run ./cmd/lunar-agent validate
```
Expected:
* Validation error is shown

### Failed shell command
```yaml
failure-test:
  repository_path: "."
  steps:
    - type: "shell"
      name: "fail"
      command: "false"
```
Expected:
* Deployment fails
* Result is saved
* Exit code is ``1``