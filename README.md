## Installation

### `cd $(echo $GOPATH)/src`

### `mkdir yourcompany.com`

### `git clone git@github.com:renelnicolas/golang-api-sample.git`

### `go run .` OR `go run main.go` OR `go build -o bin/main . && bin/main`

### build for linux `env GOOS=linux GOARCH=amd64 go build -o bin/main`

## Use api behind nginx

```
upstream backend_golang {
    server 127.0.0.1:8000;
}

server {
    listen          80;
    server_name     api.yourdomain.local;

    location / {
        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_pass http://backend_golang;
    }

    error_log /var/log/nginx/api.yourdomain.local_error.log;
    access_log /var/log/nginx/api.yourdomain.local_access.log;
}
```
