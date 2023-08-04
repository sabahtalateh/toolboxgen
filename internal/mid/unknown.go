package mid

func (t *UnknownExpr) Get() GetOnTypeRef {
	return GetOnTypeRef{typ: t}
}
