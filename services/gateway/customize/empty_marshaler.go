package customize

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type EmptyMarshaler struct {
	runtime.Marshaler
}

func (m *EmptyMarshaler) Marshal(v interface{}) ([]byte, error) {
	return nil, nil
}

func (m *EmptyMarshaler) Unmarshal(data []byte, v interface{}) error {
	return nil
}
