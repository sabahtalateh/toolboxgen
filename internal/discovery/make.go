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
		return makeComponent(calls)
	default:
		return nil, fmt.Errorf("unsupported package `%s`\n\tat %s", pkg, firstCall.position)
	}
}

func makeComponent(calls []call) (Tool, error) {
	firstCall := calls[0]
	switch firstCall.funcName {
	case "Register":
		return makeRegister(calls)
	case "RegisterE":
		return makeRegister(calls)
	case "Provider":
		return makeProvider(calls)
	case "ProviderE":
		return makeProvider(calls)
	default:
		return nil, fmt.Errorf("unsupported function `%s`\n\tat %s", firstCall.funcName, firstCall.position)
	}
}

func makeRegister(calls []call) (*component.Register, error) {
	firstCall := calls[0]
	if err := validate(
		func() error { return v.component.typed(firstCall) },
		func() error { return v.component.typeNotEmpty(firstCall) },
	); err != nil {
		return nil, err
	}

	return &component.Register{
		Type: tool.Type{
			Package: firstCall.typeParameter.pakage,
			Type:    firstCall.typeParameter.typ,
		},
		WithError: firstCall.funcName == "RegisterE",
	}, nil
}

func makeProvider(calls []call) (*component.Provider, error) {
	firstCall := calls[0]
	if err := validate(
		func() error { return v.component.typed(firstCall) },
		func() error { return v.component.typeNotEmpty(firstCall) },
	); err != nil {
		return nil, err
	}

	return &component.Provider{
		Type: tool.Type{
			Package: firstCall.typeParameter.pakage,
			Type:    firstCall.typeParameter.typ,
		},
		WithError: firstCall.funcName == "ProviderE",
	}, nil
}
