package task

type Task struct {
	Name string
	F    func(*T)
}

type T struct {
	name   string
	desc   string
	action func()
}

func (t *T) Name(name string) *T {
	t.name = name
	return t
}

func (t *T) Describe(desc string) *T {
	t.desc = desc
	return t
}

func (t *T) Action(action func()) *T {
	t.action = action
	return t
}
