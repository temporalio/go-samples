package pso

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	commonpb "go.temporal.io/temporal-proto/common"
	"go.temporal.io/temporal/encoded"
)

// jsonDataConverter implements encoded.DataConverter using JSON for Swarm and Particle
// WARGNING: Make sure all struct members are public (Capital letter) otherwise serialization does not work!
// TODO: consider storing blobs in external DB or S3
type jsonDataConverter struct {
}

// NewJSONDataConverter creates a json data converter
func NewJSONDataConverter() encoded.DataConverter {
	return &jsonDataConverter{}
}

// Json data converter implementation

func (dc *jsonDataConverter) ToData(value ...interface{}) (*commonpb.Payloads, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	var err error
	for i, obj := range value {
		switch t := obj.(type) {
		case Swarm:
			err = enc.Encode(*t.Settings)
			if err == nil {
				err = enc.Encode(*t.Gbest)
				if err == nil {
					if t.Settings.Size > 0 {
						for _, particle := range t.Particles {
							if particle == nil {
								particle = new(Particle)
							}
							err = enc.Encode(*particle)
						}
					}
				}
			}
		case WorkflowResult:
			_ = enc.Encode(t.Msg)
			err = enc.Encode(t.Success)
		default:
			err = enc.Encode(obj)
		}
		if err != nil {
			return nil, fmt.Errorf(
				"unable to encode argument: %d, %v, with error: %v", i, reflect.TypeOf(obj), err)
		}
	}

	return encoded.GetDefaultDataConverter().ToData(buf.Bytes())
	// TODO: store buf.Bytes() in DB/S3 and get key
	// return key, nil
}

func (dc *jsonDataConverter) FromData(input *commonpb.Payloads, valuePtr ...interface{}) error {
	// TODO: convert input into key in DB/S3 and retrieve bytes
	var b []byte
	err := encoded.GetDefaultDataConverter().FromData(input, &b)
	if err != nil {
		return err
	}

	if b == nil {
		return nil
	}

	dec := json.NewDecoder(bytes.NewBuffer(b))
	for i, obj := range valuePtr {
		switch t := obj.(type) {
		case *Swarm:
			t.Settings = new(SwarmSettings)
			_ = dec.Decode(t.Settings)
			t.Settings.function = FunctionFactory(t.Settings.FunctionName)
			t.Gbest = NewPosition(t.Settings.function.dim)
			err = dec.Decode(t.Gbest)
			t.Particles = make([]*Particle, t.Settings.Size)
			for index := 0; index < t.Settings.Size; index++ {
				t.Particles[index] = new(Particle)
				err = dec.Decode(t.Particles[index])
			}
		case *WorkflowResult:
			_ = dec.Decode(&t.Msg)
			err = dec.Decode(&t.Success)
		default:
			err = dec.Decode(obj)
		}
		if err != nil {
			return fmt.Errorf(
				"unable to decode argument: %d, %v, with error: %v", i, reflect.TypeOf(obj), err)
		}
	}
	return nil
}
