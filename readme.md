# simpleredirect

A basic redirect server written in go that can be used as a backend for GKE Ingress to redirect one domain to another on GCP.

## Features

- allow redirection of a domain to any url. E.g. test.com => https://example.com/someurl
- uses `/hc` to respond to health checks. This is useful for Ingress on GKE for NEGs

## Build and run

### build

```bash
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/simpleredirect .
```

### create the docker image

```bash
docker build -t gabihodoroaga/simpleredirect .
```

### run in docker 

```bash
docker run --rm -it -p 8080:8080 gabihodoroaga/simpleredirect --redirect=t1.com:https://t2.com:302
```

### deploy using kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: server-redirect
spec:
  replicas: 1
  selector:
    matchLabels:
      app: server-redirect
  strategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: server-redirect
    spec:
      containers:
      - name: server-redirect
        image: gabihodoroaga/simpleredirect:1.0.2
        args:
        - "-redirect=host.com:https://example.com"
        - "-redirect=test.com:https://example.com/test:302"
        ports:
        - name: http
          containerPort: 8080
```

### test

```bash
curl -v -H "Host: t1.com" http://localhost:8080
```

## TODO:

- [ ] add unit tests for `redirect` function
- [ ] configurable access logs

## Is it any good

I hope will be any good for someone else at least to save a half a day to build 
something similar when you don't have the time. If this is the case then you can
give me a star. It will make me smile.
