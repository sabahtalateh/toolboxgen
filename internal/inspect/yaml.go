package inspect

import (
	"bytes"
	"fmt"

	"github.com/life4/genesis/slices"
	"gopkg.in/yaml.v3"
)

type Out struct {
	Yaml *yaml.Node
}

func (o Out) String() string {
	var b bytes.Buffer
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)
	if err := enc.Encode(o.Yaml); err != nil {
		panic(err)
	}

	return b.String()
}

type node interface {
	node()
	content() []*yaml.Node
	yaml() *yaml.Node
}

type (
	scalar struct {
		value *yaml.Node
	}

	sequence struct {
		values []node
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

func (n *scalar) content() []*yaml.Node {
	return []*yaml.Node{n.value}
}

func (n *sequence) content() []*yaml.Node {
	var content []*yaml.Node
	for _, value := range n.values {
		switch value.(type) {
		case *mapping, *nesting:
			content = append(content, value.yaml())
		default:
			content = append(content, value.content()...)
		}
	}
	return []*yaml.Node{{Kind: yaml.SequenceNode, Tag: "!!seq", Content: content}}
}

func (n *mapping) content() []*yaml.Node {
	res := []*yaml.Node{{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: n.key,
	}}

	switch nn := n.value.(type) {
	case *scalar:
		res = append(res, nn.content()...)
	case *sequence:
		res = append(res, nn.content()...)
	}

	return res
}

func (n *nesting) content() []*yaml.Node {
	var content []*yaml.Node
	for _, value := range n.values {
		content = append(content, value.content()...)
	}

	res := []*yaml.Node{{Kind: yaml.ScalarNode, Tag: "!!str", Value: n.key}}

	if len(content) == 0 {
		return append(res, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"})
	}

	return append(res, &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: content})
}

func (n *scalar) yaml() *yaml.Node {
	return n.content()[0]
}

func (n *sequence) yaml() *yaml.Node {
	return n.content()[0]
}

func (n *mapping) yaml() *yaml.Node {
	return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: n.content()}
}

func (n *nesting) yaml() *yaml.Node {
	return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: n.content()}
}

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

func v(x any) node {
	return n(x)
}

func vv(xx ...any) *sequence {
	return &sequence{values: slices.Map(xx, func(el any) node { return n(el) })}
}

func n(x any) node {
	if x == nil {
		return &scalar{value: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}}
	}

	switch xx := x.(type) {
	case node:
		return xx
	case *yaml.Node:
		return &scalar{value: xx}
	case string:
		return &scalar{value: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: xx}}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return &scalar{value: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%d", xx)}}
	case float32, float64:
		return &scalar{value: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: fmt.Sprintf("%f", xx)}}
	case bool:
		b := "true"
		if !xx {
			b = "false"
		}
		return &scalar{value: &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: b}}
	default:
		panic("unknown value type")
	}
}
