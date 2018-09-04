package shared

import (
	"fmt"
	"github.com/getlantern/deepcopy"
	"github.com/golang/glog"
	"github.com/imdario/mergo"
	"sort"
	"sync"
)

type TaskFunc func(ctx *Context) error

type Tasks []*Task

type Task struct {
	Name         string
	Dependencies []string
	Func         TaskFunc

	// Contains the dependencies sorted by alphabet instead of logical.
	// Handy for comparing whether the dependencies equal.
	sortedDeps []string
}

func NewTask(name string, fun TaskFunc, deps ...string) *Task {
	sorted := make([]string, len(deps))
	copy(sorted, deps)
	sort.Strings(sorted)

	return &Task{
		Name:         name,
		Dependencies: deps,
		Func:         fun,
		sortedDeps:   sorted,
	}
}

func (tasks Tasks) GetNames() []string {
	strs := make([]string, len(tasks))

	for i, t := range tasks {
		strs[i] = t.Name
	}

	return strs
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
			go func(i2 int, t2 *Task) {
				defer wg.Done()
				taskCtx := Context{}
				deepcopy.Copy(&taskCtx, ctx)
				taskCtx.CurrentTask = t2.Name
				taskContexts[i] = taskCtx

				err := t2.Func(&taskContexts[i])
				if err != nil {
					glog.V(6).Infof("Failed executing task %s: %v", t2.Name, err)
					errors[i2] = err
					return
				}

				glog.V(6).Infof("Executed task %s", t2.Name)
			}(i, t)
		}

		wg.Wait()

		for i, taskCtx := range taskContexts {
			err := mergo.Merge(ctx, taskCtx)
			if err != nil {
				return fmt.Errorf("couldn't merge task context for task %s, see: %v", g[i].Name, err)
			}
		}

		for i, err := range errors {
			if err != nil {
				return fmt.Errorf("Error on task %s, see: %v", g[i].Name, err)
			}
		}
	}

	return nil
}

func (tasks Tasks) sort() (Tasks, error) {
	taskMap := make(map[string]*Task)

	for _, t := range tasks {
		taskMap[t.Name] = t
	}

	sorted := make(Tasks, 0, len(tasks))
	visited := make(map[string]struct{})

	for _, t := range tasks {
		err := visit(t, taskMap, visited, &sorted)
		if err != nil {
			return nil, err
		}
	}

	return sorted, nil
}

func (tasks Tasks) group() []Tasks {
	groups := make([]Tasks, 0)
	lastDeps := []string{"Cheers! 🍺"}

	for _, t := range tasks {
		if !depsEqual(lastDeps, t.sortedDeps) {
			newGroup := make(Tasks, 0)
			groups = append(groups, newGroup)
		}

		offset := len(groups) - 1
		groups[offset] = append(groups[offset], t)
		lastDeps = t.sortedDeps
	}

	return groups
}

func depsEqual(a []string, b []string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}

	if len(a) != len(b) {
		return false
	}

	for i, str := range a {
		if str != b[i] {
			return false
		}
	}

	return true
}

func visit(t *Task, all map[string]*Task, visited map[string]struct{}, sorted *Tasks) error {
	_, alreadyVisited := visited[t.Name]

	if !alreadyVisited {
		visited[t.Name] = struct{}{}

		for _, dep := range t.Dependencies {
			depT, ok := all[dep]
			if !ok {
				return fmt.Errorf("couldn't find dependency %s of task %s.", dep, t.Name)
			}

			err := visit(depT, all, visited, sorted)
			if err != nil {
				return err
			}
		}

		*sorted = append(*sorted, t)

	} else {
		found := false

		for _, t2 := range *sorted {
			if t2.Name == t.Name {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Cyclic dependency found for task %s", t.Name)
		}
	}

	return nil
}
