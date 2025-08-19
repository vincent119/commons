範例：

```go
import (
    "github.com/vincent119/commons/stringx"
    "github.com/vincent119/commons/timex"
)

func demo() {
    s := stringx.ToSnake("UserID") // "user_i_d"
    _ = s

    utc := timex.NowUTC()
    _ = utc
}
```

```
make tidy
make lint
make test
make bench

```
