package adapters

import "github.com/google/uuid"

type GuidBasedIdGenerator struct {
}

func NewGuidBasedIdGenerator() *GuidBasedIdGenerator {
	return &GuidBasedIdGenerator{}
}

func (idGen *GuidBasedIdGenerator) MakeId() string {
	return uuid.New().String()
}
