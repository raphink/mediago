Mediago
=======

Check book status for http://www.bm-chambery.fr/

Supports checking multiple accounts.

![screenshot](./screenshot.png)


# Installing

```shell
$ go get github.com/raphink/mediago
```


# Configuration example

```toml
# $HOME/.mediago.conf
renew_before = "24h"
report = "smtp"

[[account]]
name = "Foo"
login = "CHAM123456"
password = "SUPERPASS"

[[account]]
name = "Bar"
login = "CHAM456789"
password = "ANOTHERPASS"

[smtp]
username = "smtp_user"
password = "email_pass"
hostname = "smtp.example.com"
port = 587
recipients = ["foo@example.com", "bar@example.com"]
```
