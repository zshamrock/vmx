package core

const (
	AdHocCommandName = "ad-hoc"
)

type Command struct {
	Name, Command, WorkingDir string
	RequiresConfirmation      bool
}

func (c Command) IsAdHoc() bool {
	return c.Name == AdHocCommandName
}
