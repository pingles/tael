# tael

Tells a (logged) AWS ElasticSearch tale.

Provides `tail` like interface for accessing log information stored in ElasticSearch.

## Usage

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