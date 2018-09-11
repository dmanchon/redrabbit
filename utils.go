package redrabbit

import (
	"fmt"
	"github.com/satori/go.uuid"
)

func addSuffix(s, suffix string) string {
	return fmt.Sprintf("%s%s", s, suffix)
}

func GenerateUUID() string {
	u1 := uuid.Must(uuid.NewV4())
	return u1.String()
}
