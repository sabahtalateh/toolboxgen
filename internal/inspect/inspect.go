package inspect

import (
	"fmt"
	"strings"
)

type Inspect interface {
	Inspect(Context) string
}

type Context struct {
	TrimPackage string
}

func EmptyContext() Context {
	return Context{}
}

func (c Context) WithTrimPackage(p string) Context {
	return Context{TrimPackage: p}
}

func typeID(ctx Context, Package, TypeName string) string {
	if Package == ctx.TrimPackage {
		return TypeName
	}
	out := fmt.Sprintf("%s.%s", Package, TypeName)
	return strings.TrimPrefix(out, ctx.TrimPackage)
}

func typeBlock(ctx Context, Package, TypeName string) string {
	return fmt.Sprintf("type %s", typeID(ctx, Package, TypeName))
}
