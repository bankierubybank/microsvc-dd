# microsvc-dd

This project is for learning Golang (with Gin-gonic)

- How to write sample API with Golang
- How to Dockerize Golang application
- How to CI/CD

### Docker Build and Run

```
docker build -t microsvc-dd:latest .
docker run --name microsvc-dd -d -p 8080:8080 microsvc-dd:latest
```

### Docker Compose

Build and run from source

```
docker compose up -d --build
```

Or `docker compose up -d --build --force-recreate app` for rebuild only the application

Run from built image

```
docker compose up -d
```

### Landing Page

```
http://<CONTAINER-HOST-IP>
```

### Swagger

```
http://<CONTAINER-HOST-IP>:<PORT>/swagger/index.html
```

##### Environment Variable

| Variable name | Description      | Default | Mandatory |
| ------------- | ---------------- | ------- | --------- |
| PORT          | Application Port | 8080    | YES       |
| POD_NAME      | Pod Name         | unset   | NO        |
| NAMESPACE     | Namespace        | unset   | NO        |
| NODENAME      | Node Name        | unset   | NO        |
