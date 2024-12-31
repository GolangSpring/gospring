package action

type ModuleAction interface {
	Action() string
	Args() map[string]any
	AsJson() map[string]any
}
