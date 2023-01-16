package bsoncodec

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonoptions"
)

func NewStructCodec(structCodecOptions *bsonoptions.StructCodecOptions) *bsoncodec.StructCodec {
	codec, err := bsoncodec.NewStructCodec(bsoncodec.DefaultStructTagParser, structCodecOptions)
	if err != nil {
		// This function is called from the codec registration path, so errors can't be propagated. If there's an error
		// constructing the StructCodec, we panic to avoid losing it.
		panic(fmt.Errorf("error creating default StructCodec: %v", err))
	}
	return codec
}
