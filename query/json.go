package query

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/imdario/mergo"
)

var ErrInvalidData error = errors.New("query: invalid data provided")

var (
	andSplitter     *regexp.Regexp
	bracketSplitter *regexp.Regexp
)

func init() {
	andSplitter = regexp.MustCompile("&")
	bracketSplitter = regexp.MustCompile("\\[|\\]")
}

// JSON will turn a query string like:
//   cat=1&bar%5Bone%5D%5Btwo%5D=2&bar[one][red]=112
// into a JSON object with all the data merged as nicely as
// possible. Eg the example above would output:
//   {"bar":{"one":{"two":2,"red":112}}}
//
// This does not currently support arrays. Eg:
//   a[]=1&a[]=2
// will not be properly encoded into {"a":[1,2]}
func JSON(rawQuery string) ([]byte, error) {
	escapedQuery, err := url.QueryUnescape(rawQuery)
	if err != nil {
		return nil, err
	}

	builder := make(map[string]interface{})
	params := strings.Split(escapedQuery, "&")
	for _, str := range params {
		tempMap, err := queryToMap(str)
		if err != nil {
			return nil, err
		}
		err = mergo.Merge(&tempMap, builder)
		if err != nil {
			return nil, err
		}
		builder = tempMap
	}

	return json.Marshal(builder)
}

// queryToMap turns something like a[b][c]=4 into
// map[string]interface{}{
//   "a": map[string]interface{}{
// 		"b": map[string]interface{}{
// 			"c": 4,
// 		},
// 	},
// }
func queryToMap(param string) (map[string]interface{}, error) {
	temp := strings.Split(param, "=")
	if len(temp) > 2 {
		return nil, ErrInvalidData
	}

	key, valueStr := temp[0], temp[1]
	pieces := bracketSplitter.Split(key, -1)

	// nothing is nested if len == 1
	if len(pieces) == 1 {
		var value interface{}
		err := json.Unmarshal([]byte(valueStr), &value)
		if err != nil {
			// try wrapping the value in quotes
			err = json.Unmarshal([]byte("\""+valueStr+"\""), &value)
			if err != nil {
				return nil, err
			}
		}
		return map[string]interface{}{
			key: value,
		}, nil
	}
	// otherwise we have nested stuff
	ret := make(map[string]interface{}, 0)
	var err error
	ret[pieces[0]], err = queryToMap(buildNewKey(key, pieces) + "=" + valueStr)
	if err != nil {
		return nil, err
	}

	if pieces[1] == "" {
		// it is an array, eg "a[]=1"
		// This will get us the correct value for the array
		// and cleanup the key, but our merge code doesn't
		// handle arrays correctly yet.
		temp := ret[pieces[0]].(map[string]interface{})
		ret[pieces[0]] = []interface{}{temp[""]}
	}
	return ret, nil
}

// buildNewKey will take something like:
// origKey = "bar[one][two]"
// pieces = [bar one two ]
// and return "one[two]"
func buildNewKey(origKey string, pieces []string) string {
	temp := origKey[len(pieces[0])+1:]
	temp = temp[:len(pieces[1])] + temp[len(pieces[1])+1:]
	return temp
}
