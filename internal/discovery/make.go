package discovery

import (
	"fmt"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/discovery/tool"
	"github.com/sabahtalateh/toolboxgen/internal/discovery/tool/component"
)

func (d *discovery) makeTool(calls []call) (Tool, error) {
	if len(calls) == 0 {
		return nil, nil
	}

	firstCall := calls[0]
	pkg := firstCall.pakage
	if !strings.HasPrefix(pkg, "github.com/sabahtalateh/toolbox") {
		return nil, nil
	}

	parseErr := firstError(calls)
	if parseErr != nil {
		return nil, fmt.Errorf("%w\n\tat %s", parseErr.err, parseErr.position)
	}

	switch pkg {
	case "github.com/sabahtalateh/toolbox/di/component":
		return d.makeComponent(calls)
	default:
		return nil, fmt.Errorf("unsupported package `%s`\n\tat %s", pkg, firstCall.position)
	}
}

func (d *discovery) makeComponent(calls []call) (Tool, error) {
	firstCall := calls[0]
	switch firstCall.funcName {
	case "Constructor":
		return d.makeConstructor(calls)
	case "Provider":
		return d.makeProvider(calls)
	default:
		return nil, fmt.Errorf("unsupported function `%s`\n\tat %s", firstCall.funcName, firstCall.position)
	}
}

func (d *discovery) makeConstructor(calls []call) (*component.Constructor, error) {
	firstCall := calls[0]
	if err := validate(
		func() error { return v.component.typed(firstCall) },
		func() error { return v.component.typeNotEmpty(firstCall) },
		func() error { return v.component.constructorArgs(firstCall) },
		func() error { return v.component.name(calls) },
	); err != nil {
		return nil, err
	}

	// d.findFunction(firstCall.arguments[0].refArg.pakage)

	c := component.Constructor{
		Type: tool.Type{
			Package: firstCall.typeParameter.pakage,
			Type:    firstCall.typeParameter.typ,
		},
		Name: extractName(calls),
	}

	return &c, nil
}

func (d *discovery) makeProvider(calls []call) (*component.Provider, error) {
	firstCall := calls[0]
	if err := validate(
		func() error { return v.component.typed(firstCall) },
		func() error { return v.component.typeNotEmpty(firstCall) },
		func() error { return v.component.providerArgs(firstCall) },
		func() error { return v.component.name(calls) },
	); err != nil {
		return nil, err
	}

	p := component.Provider{
		Type: tool.Type{
			Package: firstCall.typeParameter.pakage,
			Type:    firstCall.typeParameter.typ,
		},
		Name: extractName(calls),
	}

	return &p, nil
}

func extractName(cc []call) string {
	for _, c := range cc {
		if c.funcName == "Name" {
			return c.arguments[0].stringArg.val
		}
	}

	return ""
}
