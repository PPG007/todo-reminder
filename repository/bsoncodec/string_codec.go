package bsoncodec

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"

	mongo_bsoncodec "go.mongodb.org/mongo-driver/bson/bsoncodec"
)

var decodeObjectIDAsHex = true

var defaultStringCodec = mongo_bsoncodec.NewStringCodec(&bsonoptions.StringCodecOptions{
	DecodeObjectIDAsHex: &decodeObjectIDAsHex,
})

type StringCodec struct {
}

func NewStringCodec() *StringCodec {
	return &StringCodec{}
}

// EncodeValue is the ValueEncoder for string types.
func (sc *StringCodec) EncodeValue(ectx mongo_bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	return defaultStringCodec.EncodeValue(ectx, vw, val)
}

// DecodeValue is the ValueDecoder for string types.
func (sc *StringCodec) DecodeValue(dctx mongo_bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	err := defaultStringCodec.DecodeValue(dctx, vr, val)
	if err == nil {
		return nil
	}

	str := ""

	switch vr.Type() {
	case bsontype.DateTime:
		dt, err := vr.ReadDateTime()
		if err != nil {
			return err
		}
		str = time.Unix(dt/1000, dt%1000*1000000).Format("2006-01-02T15:04:05.999Z07:00")
	case bsontype.Int32:
		i, err := vr.ReadInt32()
		if err != nil {
			return err
		}
		str = strconv.FormatInt(int64(i), 10)
	case bsontype.Int64:
		i, err := vr.ReadInt64()
		if err != nil {
			return err
		}
		str = strconv.FormatInt(i, 10)
	case bsontype.Double:
		i, err := vr.ReadDouble()
		if err != nil {
			return err
		}
		str = fmt.Sprintf("%f", i)
	case bsontype.Boolean:
		i, err := vr.ReadBoolean()
		if err != nil {
			return err
		}
		str = strconv.FormatBool(i)
	case bsontype.Timestamp:
		t, _, err := vr.ReadTimestamp()
		if err != nil {
			return err
		}
		str = time.Unix(int64(t), 0).Format("2006-01-02T15:04:05.999Z07:00")
	default:
		return fmt.Errorf("cannot decode %v into a string type", vr.Type())
	}

	val.SetString(str)

	return nil
}
