package inspect

import (
	"github.com/life4/genesis/slices"
	"github.com/sabahtalateh/toolboxgen/internal/types"
	"gopkg.in/yaml.v3"
)

func (i *Inspect) introField(f *types.Field) *yaml.Node {
	out := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}

	if f.Name != "" {
		// out =
	}

	return out
}

func (i *Inspect) field(f *types.Field) string {
	name := ""
	if f.Name != "" {
		name = f.Name + " "
	}
	typeOut := i.TypeRef(f.Type)
	return name + typeOut
}

func (i *Inspect) Field(f *types.Field) any {
	if i.intro {
		return i.introField(f)
	}
	return i.field(f)
}

func (i *Inspect) Field2(f *types.Field) string {
	name := ""
	if f.Name != "" {
		name = f.Name + " "
	}
	typeOut := i.TypeRef(f.Type)
	return name + typeOut
}

func (i *Inspect) Fields(f types.Fields) []any {
	return slices.Map(f, func(el *types.Field) any { return i.Field(el) })
}
