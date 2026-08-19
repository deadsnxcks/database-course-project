package dto

import (
	"time"

	operationv1 "github.com/deadsnxcks/dbcp/protos/gen/go/operation"
)

type Operation struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
}

func OperationFromProto(o *operationv1.Operation) Operation {
	return Operation{
		ID:        o.GetId(),
		Title:     o.GetTitle(),
		CreatedAt: o.GetCreatedAt().AsTime(),
	}
}

func OperationsFromProto(ops []*operationv1.Operation) []Operation {
	out := make([]Operation, 0, len(ops))
	for _, o := range ops {
		out = append(out, OperationFromProto(o))
	}

	return out
}