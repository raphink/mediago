Mediago
=======

Check book status for http://www.bm-chambery.fr/

Supports checking multiple accounts.


# Installing

```shell
$ go get github.com/raphink/mediago
```


# Configuration example

```toml
# $HOME/.mediago.conf
[[account]]
name = "Foo"
login = "CHAM123456"
password = "SUPERPASS"

[[account]]
name = "Bar"
login = "CHAM456789"
password = "ANOTHERPASS"
```
