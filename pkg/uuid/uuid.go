package uuid

import (
	"github.com/google/uuid"
)

func NewUID() string {
	return uuid.New().String()
}
