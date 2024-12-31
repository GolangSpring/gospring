package repository

import (
	"errors"
	"go-spring/domains/action"
)

type Playbook struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Become bool   `json:"become"`
	Plays  []Play `json:"plays"`
}

type Play struct {
	Name           string         `json:"name"`             // Optional name of the play
	Hosts          string         `json:"hosts"`            // Target hosts
	Tasks          []*Task        `json:"tasks"`            // Tasks to execute
	Variables      map[string]any `json:"vars"`             // Variables for the play
	AnyErrorsFatal bool           `json:"any_errors_fatal"` // Whether any error should stop execution
}

type Task struct {
	Name         string              `json:"name"`          // Optional name of the task
	ModuleAction action.ModuleAction `json:"module_action"` // Module action to execute
	When         string              `json:"when"`          // Condition for execution
	WithItems    []any               `json:"with_items"`    // Items for looping
}

func (play *Play) InsertTask(insertedAt int, task *Task) error {
	if insertedAt < 0 || insertedAt > len(play.Tasks)-1 {
		return errors.New("index out of range")
	}

	newTasks := []*Task{}
	for idx, currentTask := range play.Tasks {
		if idx == insertedAt {
			newTasks = append(newTasks, task)
		} else {
			newTasks = append(newTasks, currentTask)
		}
	}
	play.Tasks = newTasks

	return nil
}

func (play *Play) RemoveTask(removeAt int) error {
	if removeAt < 0 || removeAt > len(play.Tasks)-1 {
		return errors.New("index out of range")
	}
	newTasks := []*Task{}
	for idx, currentTask := range play.Tasks {
		if idx == removeAt {
			continue
		}
		newTasks = append(newTasks, currentTask)
	}
	play.Tasks = newTasks
	return nil
}

func (play *Play) AppendTask(task *Task) {
	play.Tasks = append(play.Tasks, task)
}

func (play *Play) PopTask() {
	play.Tasks = play.Tasks[1:]
}

func (play *Play) AsJson() map[string]any {
	playMap := map[string]any{
		"name":             play.Name,
		"hosts":            play.Hosts,
		"vars":             play.Variables,
		"any_errors_fatal": play.AnyErrorsFatal,
	}

	// Convert tasks to their JSON representation
	taskMaps := []map[string]any{}
	for _, task := range play.Tasks {
		taskMaps = append(taskMaps, task.AsJson())
	}
	playMap["tasks"] = taskMaps

	return playMap
}

func (book *Playbook) AsJson() []map[string]any {

	playMaps := []map[string]any{}
	for _, play := range book.Plays {
		playMaps = append(playMaps, play.AsJson())
	}
	return playMaps
}

func (task *Task) AsJson() map[string]any {
	taskMap := map[string]any{
		"name": task.Name,
	}

	if len(task.When) > 0 {
		taskMap["when"] = task.When
	}

	if len(task.WithItems) > 0 {
		taskMap["with_items"] = task.WithItems
	}

	// Flatten module_action into the task map
	actionMap := task.ModuleAction.AsJson()
	for key, value := range actionMap {
		taskMap[key] = value
	}
	return taskMap
}
