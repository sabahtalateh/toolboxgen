package discovery

import (
	"fmt"
	"github.com/sabahtalateh/toolboxgen/internal/discovery/tool"
	"github.com/sabahtalateh/toolboxgen/internal/discovery/tool/component"
	"strings"
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
		return nil, fmt.Errorf("%w\n\tat%s", parseErr.err, parseErr.position)
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
	default:
		return nil, fmt.Errorf("unsupported function `%s`\n\tat %s", firstCall.funcName, firstCall.position)
	}
}

func makeRegister(calls []call) (Tool, error) {
	firstCall := calls[0]
	if err := validate(
		func() error { return v.component.register.typed(firstCall) },
		func() error { return v.component.register.typeNotEmpty(firstCall) },
	); err != nil {
		return nil, err
	}

	return &component.Register{
		Type: tool.Type{
			Package: firstCall.typeParameter.pakage,
			Type:    firstCall.typeParameter.typ,
		},
	}, nil
}
