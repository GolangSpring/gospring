package action

type CommandModuleAction struct {
	Command string `json:"command"`
}

func (action *CommandModuleAction) Action() string {
	return "ansible.builtin.command"
}

func (action *CommandModuleAction) Args() map[string]any {
	return map[string]any{
		"cmd": action.Command,
	}
}

func (action *CommandModuleAction) AsJson() map[string]any {
	return map[string]any{
		action.Action(): action.Args(),
	}
}
