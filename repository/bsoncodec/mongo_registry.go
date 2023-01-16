package bsoncodec

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"
)

var DefaultRegistry *bsoncodec.Registry = nil

func init() {
	tTime := reflect.TypeOf(time.Time{})

	rb := bsoncodec.NewRegistryBuilder()
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)

	//decoder
	useLocalTimeZone := true
	decodeZeroStruct := true
	rb.RegisterDefaultDecoder(reflect.String, NewStringCodec())
	rb.RegisterDefaultDecoder(reflect.Struct, NewStructCodec(&bsonoptions.StructCodecOptions{DecodeZeroStruct: &decodeZeroStruct}))
	rb.RegisterTypeDecoder(tTime, bsoncodec.NewTimeCodec(&bsonoptions.TimeCodecOptions{UseLocalTimeZone: &useLocalTimeZone}))
	rb.RegisterTypeDecoder(reflect.TypeOf(ObjectId("")), bsoncodec.ValueDecoderFunc(ObjectIDDecodeValue))

	//encoder
	encodeNilAsEmpty := true
	encodeOmitDefaultStruct := true
	rb.RegisterDefaultEncoder(reflect.Struct, NewStructCodec(&bsonoptions.StructCodecOptions{EncodeOmitDefaultStruct: &encodeOmitDefaultStruct}))
	rb.RegisterDefaultEncoder(reflect.Slice, bsoncodec.NewSliceCodec(&bsonoptions.SliceCodecOptions{EncodeNilAsEmpty: &encodeNilAsEmpty}))
	rb.RegisterTypeEncoder(reflect.TypeOf(ObjectId("")), bsoncodec.ValueEncoderFunc(ObjectIDEncodeValue))

	//map
	rb.RegisterTypeMapEntry(bsontype.Array, reflect.TypeOf([]interface{}(nil)))
	rb.RegisterTypeMapEntry(bsontype.DateTime, tTime)
	rb.RegisterTypeMapEntry(bsontype.Int32, reflect.TypeOf(int(0)))
	rb.RegisterTypeMapEntry(bsontype.EmbeddedDocument, reflect.TypeOf(bson.M{}))
	rb.RegisterTypeMapEntry(bsontype.ObjectID, reflect.TypeOf(ObjectId("")))

	var primitiveCodecs bson.PrimitiveCodecs

	primitiveCodecs.RegisterPrimitiveCodecs(rb)

	DefaultRegistry = rb.Build()
}

func ObjectIDEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	oid := primitive.NilObjectID

	if val.Interface().(ObjectId).Hex() == "" {
		return bsoncodec.ValueEncoderError{Name: "ObjectIDEncodeValueWithZero", Kinds: []reflect.Kind{reflect.ValueOf(ObjectId("")).Kind()}, Received: val}
	}

	if val.Interface().(ObjectId).Valid() {
		var err error
		oid, err = primitive.ObjectIDFromHex(val.Interface().(ObjectId).Hex())
		if err != nil {
			return bsoncodec.ValueEncoderError{Name: fmt.Sprintf("ObjectIDEncodeValue error %v", err), Types: []reflect.Type{reflect.TypeOf(ObjectId(""))}, Received: val}
		}
	}

	return vw.WriteObjectID(oid)
}

func ObjectIDDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != reflect.TypeOf(ObjectId("")) {
		return bsoncodec.ValueDecoderError{Name: "ObjectIDDecodeValue", Types: []reflect.Type{reflect.TypeOf(ObjectId(""))}, Received: val}
	}

	var oid primitive.ObjectID
	var err error
	switch vrType := vr.Type(); vrType {
	case bsontype.ObjectID:
		oid, err = vr.ReadObjectID()
		if err != nil {
			return err
		}
	case bsontype.String:
		str, err := vr.ReadString()
		if err != nil {
			return err
		}

		if len(str) == 12 {
			byteArr := []byte(str)
			copy(oid[:], byteArr)
		} else if len(str) > 0 {
			oid, err = primitive.ObjectIDFromHex(str)
			if err != nil {
				return err
			}
		}
	case bsontype.Null:
		if err = vr.ReadNull(); err != nil {
			return err
		}
	case bsontype.Undefined:
		if err = vr.ReadUndefined(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("cannot decode %v into an ObjectID", vrType)
	}

	if oid.IsZero() {
		val.Set(reflect.ValueOf(ObjectId("")))
	} else {
		val.Set(reflect.ValueOf(ObjectIdHex(oid.Hex())))
	}

	return nil
}
