# tael

Tells a (logged) AWS ElasticSearch tale.

Provides `tail` like interface for accessing log information stored in ElasticSearch.

## Usage

```sh
$ tael -f ioexception
```

`ioexception` could be any valid ElasticSearch query. Output can be controlled by specifying:

* `-j` Outputs log entries as JSON (potentially to pipe into `jq` or other commands)
* `-l` Allows specification of the output format using Go's mustach-like templating. For example, to output just the entry IDs you could use: `tael -f -l "{{.Id}}" ...`

### Full usage:

```
usage: tael [<flags>] [<filter>...]

Flags:
      --help           Show context-sensitive help (also try --help-long and --help-man).
      --host=HOST      aws elasticsearch url
  -f, --follow         follow log
      --index="*"      elasticsearch index
  -n, --number=10      number of results to retrieve
  -j, --json           output as json
  -l, --layout=LAYOUT  custom templated output
      --query="*"      elasticsearch query

Args:
  [<filter>]  field filter, name=value
```

## Configuring

`tael` expects your ElasticSearch indices contain the following fields:

* `@timestamp`: the time of the log entry
* `message`: the log message
* `level_name`: log level, e.g. info, warn, error etc.

Additionally, `tael` expects log messages from Docker containers and so will also extract:

* `image_name`
* `container_name`


## License

Please see [LICENSE](./LICENSE)

## Contributors

* [Paul Ingles](https://github.com/pingles)
* [Siddharth Dawara](https://github.com/sdawara)
