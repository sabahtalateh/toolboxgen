package inspect

import (
	"fmt"
	"strings"
)

type Config struct {
	TrimPackage string
}

type Inspect struct {
	trimPackage string
}

func New(c Config) *Inspect {
	return &Inspect{trimPackage: c.TrimPackage}
}

func (i *Inspect) typeID(Package, TypeName string) string {
	if Package == i.trimPackage {
		return TypeName
	}
	out := fmt.Sprintf("%s.%s", Package, TypeName)
	return strings.TrimPrefix(out, i.trimPackage)
}
