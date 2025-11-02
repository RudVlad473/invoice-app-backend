package pkg

type Command interface {
	Execute()
	Undo()
}
