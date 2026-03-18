package core

import (
	"github.com/bono/loadstar/internal"
	"github.com/bono/loadstar/internal/address"
)

// Element represents a single LOADSTAR element (Map, WayPoint, Link, SavePoint).
type Element struct {
	Type    string // M, W, L, S
	ID      string
	Address string // e.g. W://root/dev/auth
	Status  string // S_IDL, S_PRG, S_STB, S_ERR, S_REV
}

// ElementService handles element lifecycle operations.
type ElementService struct {
	storage internal.Storage
	parser  func(raw string) (*address.Address, error)
}

func NewElementService(storage internal.Storage) *ElementService {
	return &ElementService{
		storage: storage,
		parser:  address.Parse,
	}
}

func (s *ElementService) Create(elementType, id, parent string) error {
	// Implemented in cmd/element.go — ElementService provides storage access.
	return nil
}

func (s *ElementService) Edit(addr string) error {
	// Implemented in cmd/element.go — ElementService provides storage access.
	return nil
}

func (s *ElementService) Delete(addr string) error {
	// Implemented in cmd/element.go — ElementService provides storage access.
	return nil
}

// Storage exposes the injected storage for use by cmd layer.
func (s *ElementService) Storage() internal.Storage {
	return s.storage
}

// ParseAddress is a convenience wrapper around the injected parser.
func (s *ElementService) ParseAddress(raw string) (*address.Address, error) {
	return s.parser(raw)
}
