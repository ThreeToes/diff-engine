diff-engine
===========
A small app to help track stealth edits from news outlets

Building
--------
Run `dep ensure` followed by `go build diff-engine.go`

Running
--------
`diff-engine` requires a JSON config file to work. This is passed in using the `-config` option with a path to the file

The file itself has the structure
```
{
  "database": {
    "host": "dbhost",
    "port": 5432,
    "name": "dbname",
    "username": "username",
    "password": "password",
    "ssl_mode": "disable"
  },
  "feeds" : {
    "vice": {
        "url":"https://www.vice.com/en_us/rss",
        "article_selector": ".article__body"
    }
  }
}
```
Where the feeds field is a key-value pair dictionary