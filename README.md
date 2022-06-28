# Simple Reverse Proxy
# Usage
Quick start
```command
$ go run reverse_proxy.go
```
Run the following command for argument explanations.
```command
$ go run reverse_proxy.go -h
```
Besides, keys of configuration file adopt two rules to locate target rules.
First, search for explicit rules: host/path/in/url.
Then, search for host: host.
If none of them are found, use default routing rule.
