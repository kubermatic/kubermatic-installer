package shared

import (
	"fmt"
	"github.com/getlantern/deepcopy"
	"github.com/golang/glog"
	"github.com/imdario/mergo"
	"sync"
)

type Tasks []Task

type BaseTask struct {
	Dependencies Tasks
}

func (t *BaseTask) GetDependencies() Tasks {
	return t.Dependencies
}

func (t *BaseTask) SetDependencies(tasks Tasks) {
	t.Dependencies = tasks
}

type Task interface {
	GetDependencies() Tasks
	SetDependencies(Tasks)

	Execute(ctx *Context) error
}

func (tasks Tasks) DumpGroups() ([]Tasks, error) {
	sorted, err := tasks.sort()
	if err != nil {
		return nil, fmt.Errorf("couldn't sort tasks, see: %v", err)
	}

	grouped := sorted.group()

	return grouped, nil
}

func (tasks Tasks) Execute(ctx *Context) error {
	sorted, err := tasks.sort()
	if err != nil {
		return fmt.Errorf("couldn't sort tasks, see: %v", err)
	}

	grouped := sorted.group()

	for _, g := range grouped {
		var wg sync.WaitGroup

		wg.Add(len(g))
		errors := make([]error, len(g))
		taskContexts := make([]Context, len(g))

		for i, t := range g {
			go func(i2 int, t2 Task) {
				defer wg.Done()
				taskCtx := Context{}
				deepcopy.Copy(&taskCtx, ctx)
				taskContexts[i] = taskCtx

				err := t2.Execute(&taskContexts[i])
				if err != nil {
					glog.V(6).Infof("Failed executing task %T: %v", t2, err)
					errors[i2] = err
					return
				}

				glog.V(6).Infof("Executed task %T", t2)
			}(i, t)
		}

		wg.Wait()

		for i, taskCtx := range taskContexts {
			err := mergo.Merge(ctx, taskCtx)
			if err != nil {
				return fmt.Errorf("couldn't merge task context for task %T, see: %v", g[i], err)
			}
		}

		for i, err := range errors {
			if err != nil {
				return fmt.Errorf("Error on task %T, see: %v", g[i], err)
			}
		}
	}

	return nil
}

func (tasks Tasks) sort() (Tasks, error) {
	sorted := make(Tasks, 0, len(tasks))
	visited := make(map[interface{}]struct{})

	for _, t := range tasks {
		err := visit(t, visited, &sorted)
		if err != nil {
			return nil, err
		}
	}

	return sorted, nil
}

func (tasks Tasks) group() []Tasks {
	groups := make([]Tasks, 0)
	lastDeps := []Task{nil}

	for _, t := range tasks {
		if !depsEqual(lastDeps, t.GetDependencies()) {
			newGroup := make(Tasks, 0)
			groups = append(groups, newGroup)
		}

		offset := len(groups) - 1
		groups[offset] = append(groups[offset], t)
		lastDeps = t.GetDependencies()
	}

	return groups
}

func depsEqual(a []Task, b []Task) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}

	if len(a) != len(b) {
		return false
	}

	for _, t1 := range a {
		found := false

		for _, t2 := range b {
			if t2 == t1 {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func visit(t Task, visited map[interface{}]struct{}, sorted *Tasks) error {
	_, alreadyVisited := visited[t]

	if !alreadyVisited {
		visited[t] = struct{}{}

		for _, depT := range t.GetDependencies() {
			err := visit(depT, visited, sorted)
			if err != nil {
				return err
			}
		}

		*sorted = append(*sorted, t)

	} else {
		found := false

		for _, t2 := range *sorted {
			if t2 == t {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Cyclic dependency found for task %T", t)
		}
	}

	return nil
}
