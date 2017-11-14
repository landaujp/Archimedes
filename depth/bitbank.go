package depth

import (
	"strconv"

	simplejson "github.com/bitly/go-simplejson"
)

type Bitbank struct {
	Json *simplejson.Json
}

func (b *Bitbank) GetDepth() int64 {
	a := b.Json.Get("data").Get("last").MustString()
	res, _ := strconv.ParseInt(a, 10, 64)
	return res
}

func (b *Bitbank) GetTimestamp() int64 {
	a := b.Json.Get("data").Get("timestamp").MustInt64()
	t := a / 1000
	return t
}

func (b *Bitbank) SetJson(json *simplejson.Json) {
	b.Json = json
}
