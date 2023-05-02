package discovery

func (d *discovery) typesFromCalls() {
	for _, c := range d.calls {
		typeFromCalls(c)
	}

}

func typeFromCalls(c callSequence) {
	switch c.pakage {
	case "github.com/sabahtalateh/toolbox/di/component":
		componentFromCalls(c)
	case "github.com/sabahtalateh/toolbox/http/route":

	}
}

func componentFromCalls(c callSequence) {
	println(&c)
}
