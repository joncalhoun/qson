# query-to-string

Convert query strings into JSON strings so that you can more easily parse them into structs/maps in Go.

I wrote this to help someone in the Gopher Slack, so it isn't really 100% complete but is meant as a starting point. It does not support arrays very well. eg:

```
a[]=1&a[]=2
```

Will **NOT** result in `{"a":[1,2]}` but will instead probably just result in `{"a":[2]}`. The merge code needs worked on for this use case.

## Usage

```go
import "github.com/joncalhoun/query-to-json/query"

// ...

b, err := query.JSON("bar%5Bone%5D%5Btwo%5D=2&bar[one][red]=112")
if err != nil {
  panic(err)
}
fmt.Println(string(b))
// Should output: {"bar":{"one":{"red":112,"two":2}}}
```
