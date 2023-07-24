package inspect

import (
	"fmt"
	"strings"
)

type Config struct {
	TrimPackage   string
	Introspective bool
}

type Inspect struct {
	trimPackage string
	intro       bool
}

func New(c Config) *Inspect {
	return &Inspect{trimPackage: c.TrimPackage, intro: c.Introspective}
}

func (i *Inspect) typeID(Package, TypeName string) string {
	if Package == i.trimPackage {
		return TypeName
	}
	out := fmt.Sprintf("%s.%s", Package, TypeName)
	return strings.TrimPrefix(out, i.trimPackage)
}

func (i *Inspect) typeBlock(Package, TypeName string) string {
	return fmt.Sprintf("type %s", i.typeID(Package, TypeName))
}
