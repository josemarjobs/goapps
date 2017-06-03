## JWT with Golang


```
$ go get gopkg.in/square/go-jose.v2
```

```
$ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem
```

Generate an RSA private key

```
$ openssl genrsa -out demo.rsa 1024
```

Generate the public key for the private key

```
$ openssl rsa -in demo.rsa -pubout > demo.rsa.pub
```