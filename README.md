socks5map
=========

Development SOCKS5 server to handle fake DNS resolution. Thin layer over https://github.com/armon/go-socks5

0. Start `socks5map`, by default listens to `:8080`, can be controller with `-l` flag.
0. Configure browser to use SOCKS5 proxy (ie. `127.0.0.1:8080`).
0. Now any hostname will be "proxied" to localhost, where you should have a web server.
0. You can also resolve particular domains to particular IPs instead of a wildcard using the `-r` flag.

Usage
-----

```
socks5map [-h] [-l addr] [-r rules]
```

* `-h` Help screen
* `-l` Listen address (default ":8080")
* `-r` Comma separated list of "domain:IP" for DNS resolving; `*` or empty matches any name. Domains not matched by any rule resolve through regular system DNS. (default `:127.0.0.1`, meaning anything to localhost)

Example:
--------

Listen to 127.0.0.2:9000, resolve hostname *home* to IP `127.0.0.1` and *extension.domain.com* to `192.168.0.100`, other domains will be resolved by system's DNS.

```
socks5map -l 127.0.0.2:9000 -r 'home:127.0.0.1,extension.domain.com:192.168.0.100'
```

Installation
-----

```
go get github.com/armon/go-socks5
go get github.com/dberstein/socks5map
```
