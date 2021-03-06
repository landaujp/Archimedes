package last

import simplejson "github.com/bitly/go-simplejson"

type Coincheck struct {
	Json *simplejson.Json
}

func (c *Coincheck) GetLast() int64 {
	a := c.Json.Get("last").MustFloat64()
	return int64(a)
}

func (c *Coincheck) GetTimestamp() int64 {
	a := c.Json.Get("timestamp").MustInt64()
	return a
}

func (c *Coincheck) SetJson(json *simplejson.Json) {
	c.Json = json
}
