package yamled

type Step interface{}

type Path []Step

func (p Path) Parent() Path {
	if len(p) < 1 {
		return nil
	}

	return p[0 : len(p)-1]
}

func (p Path) Tail() Step {
	if len(p) == 0 {
		return nil
	}

	return p[len(p)-1]
}
