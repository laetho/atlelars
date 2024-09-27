# Go webserver

Basic http server with a muxer and static fileserver.


## Usage

```bash
go run main.go
```

Get static index content:
```bash
curl http://localhost:8080/
```

Post a json payload:
```bash
curl -X POST http://localhost:8080/action/test \
  -H "Content-Type: application/json" \
  -d '{"action": "hit me!"}'
```

## Static Build

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o basic 
```

This produces a 8MB binary with only the following os dependencies on Linux:


```
❯ ldd basic
        linux-vdso.so.1 (0x0000738f0aaf1000)
        libresolv.so.2 => /usr/lib/libresolv.so.2 (0x0000738f0aaa4000)
        libc.so.6 => /usr/lib/libc.so.6 (0x0000738f0a8b3000)
        /lib64/ld-linux-x86-64.so.2 => /usr/lib64/ld-linux-x86-64.so.2 (0x0000738f0aaf3000)

```
# atlelars
