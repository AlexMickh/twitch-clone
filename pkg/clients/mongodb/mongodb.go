package mongodb

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/AlexMickh/twitch-clone/pkg/utils/retry"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func New(ctx context.Context, host string, port int, user string, password string) (*mongo.Client, error) {
	const op = "storage.mongo.New"

	var client *mongo.Client
	connString := fmt.Sprintf("mongodb://%s:%s@%s:%d/?authSource=admin", user, password, host, port)

	err := retry.WithDelay(5, 500*time.Millisecond, func() error {
		var err error

		client, err = mongo.Connect(options.Client().ApplyURI(connString).SetRegistry(UUIDRegistry))
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return client, nil
}

var (
	tUUID        = reflect.TypeOf(uuid.UUID{})
	uuidSubtype  = byte(0x04)
	UUIDRegistry = bson.NewRegistry()
)

func init() {
	UUIDRegistry.RegisterTypeEncoder(tUUID, bson.ValueEncoderFunc(uuidEncodeValue))
	UUIDRegistry.RegisterTypeDecoder(tUUID, bson.ValueDecoderFunc(uuidDecodeValue))
}

func uuidEncodeValue(ec bson.EncodeContext, vw bson.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tUUID {
		return bson.ValueEncoderError{Name: "uuidEncodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}
	b := val.Interface().(uuid.UUID)
	return vw.WriteBinaryWithSubtype(b[:], uuidSubtype)
}

func uuidDecodeValue(dc bson.DecodeContext, vr bson.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tUUID {
		return bson.ValueDecoderError{Name: "uuidDecodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}

	var data []byte
	var subtype byte
	var err error
	switch vrType := vr.Type(); vrType {
	case bson.TypeBinary:
		data, subtype, err = vr.ReadBinary()
		if subtype != uuidSubtype {
			return fmt.Errorf("unsupported binary subtype %v for UUID", subtype)
		}
	case bson.TypeNull:
		err = vr.ReadNull()
	case bson.TypeUndefined:
		err = vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a UUID", vrType)
	}

	if err != nil {
		return err
	}
	uuid2, err := uuid.FromBytes(data)
	if err != nil {
		return err
	}
	val.Set(reflect.ValueOf(uuid2))
	return nil
}
