package db_code

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const name = "json"

var Codec = new(jsonCodec)

type jsonCodec int

// Marshal Encode value with protocol buffer.
// If type isn't a Protocol buffer Message, json encoder will be used instead.
func (c jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c jsonCodec) Unmarshal(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}

func (c jsonCodec) Name() string {
	return name
}
