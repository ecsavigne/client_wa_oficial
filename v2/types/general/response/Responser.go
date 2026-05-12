package response

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type KernelResponser struct {
	xxx_type string
	response proto.Message
}

type Responser interface {
	String() string
	GetType() string
	SetType(typeStr string)
	GetResponse() proto.Message
}

func NewResponse(response proto.Message, typeStr string) Responser {
	return &KernelResponser{response: response, xxx_type: typeStr}
}

func (k *KernelResponser) GetType() string { return k.xxx_type }

// func (k *KernelResponser) String() string { return protojson.Format(k.response) }
func (k *KernelResponser) String() string {
	temp, mapResponse := make(map[string]any), make(map[string]any)

	mapResponse["type"] = k.GetType()
	formatByte, _ := protojson.MarshalOptions{UseProtoNames: true}.Marshal(k.response)

	json.Unmarshal(formatByte, &temp)
	mapResponse["meta_content"] = temp

	data, _ := json.MarshalIndent(mapResponse, "", " ")
	return string(data)
}

func (k *KernelResponser) SetType(typeStr string) { k.xxx_type = typeStr }

// func (k *KernelResponser) GetResponse() proto.Message { return k.response }
func (k *KernelResponser) GetResponse() proto.Message { return k.response }
