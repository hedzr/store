# test for consul

## The Way

### Install and run consul service

For macOS, using homebrew is simplest.

```bash
brew install consul
brew service start consul
# ...
brew service stop consul
```

### Run me

```bash
consul kv put 'ops/config/common' '---'
```

```bash
go run .
# Or: go run . -watch
```

### Insertions on consul k/v store

In another terminal window, making some insertions:

```bash
# consul kv put 'testconsul' 'debug: false'
consul kv put 'testconsul/whoami' 'name: store'
consul kv put 'testconsul/versions' 'version: v1'
```

You would saw some outputs like this:

![image-20240207191408946](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20240207191408946.png)

And try mofifying and deleting now:

```bash
consul kv put 'testconsul/please' 'wake-me: on'
consul kv put 'testconsul/please' 'wake-me: off'
consul kv delete 'testconsul/please'
```

![image-20240207191615307](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20240207191615307.png)

Using watching
