package inspect

import (
	"fmt"

	"github.com/life4/genesis/slices"
	"gopkg.in/yaml.v3"
)

type node interface {
	node()
}

type (
	scalar struct {
		n *yaml.Node
	}

	sequence struct {
		nn []*yaml.Node
	}

	mapping struct {
		key   string
		value node
	}

	nesting struct {
		key    string
		values []node
	}
)

func (n *scalar) node()   {}
func (n *sequence) node() {}
func (n *mapping) node()  {}
func (n *nesting) node()  {}

func kv(k string, vv ...node) node {
	for _, x := range vv {
		switch xx := x.(type) {
		case *scalar:
			if len(vv) > 1 {
				panic("multiple scalars")
			}
			return &mapping{key: k, value: xx}
		case *sequence:
			if len(vv) > 1 {
				panic("multiple sequences")
			}
			return &mapping{key: k, value: xx}
		}
	}

	return &nesting{key: k, values: vv}
}

func v(x any) *scalar {
	if x == nil {
		return &scalar{n: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}}
	}
	return &scalar{n: n(x)}
}

func vv(xx ...any) *sequence {
	return &sequence{nn: slices.Map(xx, func(el any) *yaml.Node { return n(el) })}
}

func n(x any) *yaml.Node {
	switch xx := x.(type) {
	case string:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: xx}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%d", xx)}
	case float32, float64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: fmt.Sprintf("%f", xx)}
	case bool:
		b := "true"
		if !xx {
			b = "false"
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: b}
	case *yaml.Node:
		return xx
	default:
		panic("unknown value type")
	}
}
