package model

import (
	"movieexample.com/gen"
	model "movieexample.com/metadata/pkg"
)

// MetadataToProto converts a Metadata struct into a generated proto counterpart.
func MetadataToProto(m *model.Metadata) *gen.Metadata {
	return &gen.Metadata{
		Id:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Director:    m.Director,
	}
}

// MetadataFromProto converts a generated proto counterpart
// into a Metadata struct.
func MetadataFromProto(m *gen.Metadata) *model.Metadata {
	return &model.Metadata{
		ID:          m.Id,
		Title:       m.Title,
		Description: m.Description,
		Director:    m.Director,
	}
}
