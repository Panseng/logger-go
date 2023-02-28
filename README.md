# usage

```shell
# get package
go get github.com/Panseng/logger-go

# import
import "github.com/Panseng/logger-go"
```

Use in code
```go
var log = logger.GetLogger("test")
log.Debugf("test")
```

Run
```shell
# data file's base dir path
export TEST_PATH="./data"

# log config
export TEST_LOG_LEVEL=debug
export TEST_LOG_FOLDER=logs
export TEST_LOG_MAXSIZE=31 # 31 Mb
export TEST_LOG_MAXBACKUPS=20
export TEST_LOG_MAXAGE=30
export TEST_LOG_LOCALTIME=true
export TEST_LOG_COMPRESS=true
export TEST_LOG_LOGINCONSOLE=true # whether print in console 是否打印到控制面板

go run .
```