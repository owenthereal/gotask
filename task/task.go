package task

type T struct {
	desc string
	name string
}

func (t *T) Describe(desc string) *T {
	t.desc = desc
	return t
}

func (t *T) Name(name string) *T {
	t.name = name
	return t
}
