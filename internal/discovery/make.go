package discovery

import (
	"fmt"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/discovery/tool"
	"github.com/sabahtalateh/toolboxgen/internal/discovery/tool/component"
)

func makeTool(calls []call) (Tool, error) {
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
		return makeComponentTool(calls)
	default:
		return nil, fmt.Errorf("unsupported package `%s`\n\tat %s", pkg, firstCall.position)
	}
}

func makeComponentTool(calls []call) (Tool, error) {
	firstCall := calls[0]
	switch firstCall.funcName {
	case "Component":
		return makeComponent(calls)
	case "Provider":
		return makeProvider(calls)
	default:
		return nil, fmt.Errorf("unsupported function `%s`\n\tat %s", firstCall.funcName, firstCall.position)
	}
}

func makeComponent(calls []call) (*component.Component, error) {
	firstCall := calls[0]
	if err := validate(
		func() error { return v.component.typed(firstCall) },
		func() error { return v.component.typeNotEmpty(firstCall) },
		func() error { return v.component.componentCall(firstCall) },
		func() error { return v.component.name(calls) },
	); err != nil {
		return nil, err
	}

	c := component.Component{
		Type: tool.Type{
			Package: firstCall.typeParameter.pakage,
			Type:    firstCall.typeParameter.typ,
		},
		Name: extractName(calls),
	}

	return &c, nil
}

func makeProvider(calls []call) (*component.Provider, error) {
	firstCall := calls[0]
	if err := validate(
		func() error { return v.component.typed(firstCall) },
		func() error { return v.component.typeNotEmpty(firstCall) },
		func() error { return v.component.providerCall(firstCall) },
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
