package task

type Task struct {
	Name        string
	Usage       string
	Description string
	F           func(*T)
}

type T struct {
	Err string
}

func (t *T) Error(err string) {
	t.Err = err
}
