package command

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/kubermatic/kubermatic-installer/pkg/shared"
	"github.com/kubermatic/kubermatic-installer/pkg/tasks"
	"gopkg.in/yaml.v2"
	"strings"
)

var tasksToExecute = make(shared.Tasks, 0)

func InstallCommand(manifestContent []byte) error {
	manifest := &shared.Manifest{}
	err := yaml.Unmarshal(manifestContent, manifest)
	if err != nil {
		return fmt.Errorf("Couldn't parse manifest, see: %v.", err)
	}

	taskCtx := &shared.Context{
		Manifest: manifest,
	}

	setupTasks(taskCtx)
	printTaskExecutionFlow()

	err = tasksToExecute.Execute(taskCtx)
	if err != nil {
		return err
	}

	return nil
}

func addTask(t shared.Task, deps ...shared.Task) shared.Task {
	t.SetDependencies(deps)
	tasksToExecute = append(tasksToExecute, t)
	return t
}

func printTaskExecutionFlow() {
	builder := &strings.Builder{}
	dumped, _ := tasksToExecute.DumpGroups()

	for _, g := range dumped {
		builder.WriteString("[\n")

		for _, t := range g {
			builder.WriteString(fmt.Sprintf("\t%#v\n", t))
		}

		builder.WriteString("]\n")
	}

	glog.V(6).Infof("Executing Tasks:\n%s", builder.String())
}

func setupTasks(ctx *shared.Context) {
	validateVersion := addTask(&tasks.ValidateVersionTask{})

	// Example, delete this when we have  more.
	addTask(&tasks.ValidateVersionTask{}, validateVersion)
	addTask(&tasks.ValidateVersionTask{}, validateVersion)
}
