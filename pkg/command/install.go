package command

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/kubermatic/kubermatic-installer/pkg/shared"
	"github.com/kubermatic/kubermatic-installer/pkg/tasks"
	"gopkg.in/yaml.v2"
	"strings"
)

func InstallCommand(manifestContent []byte) error {
	manifest := &shared.Manifest{}
	err := yaml.Unmarshal(manifestContent, manifest)
	if err != nil {
		return fmt.Errorf("Couldn't parse manifest, see: %v.", err)
	}

	taskCtx := &shared.Context{
		Manifest: manifest,
	}

	t := setupTasks(taskCtx)
	printTaskExecutionFlow(t)

	err = t.Execute(taskCtx)
	if err != nil {
		return err
	}

	return nil
}

func setupTasks(ctx *shared.Context) shared.Tasks {
	return shared.Tasks{
		shared.NewTask(tasks.TASK_VALIDATE_VERSION, tasks.ValidateVersion),
	}
}

func printTaskExecutionFlow(tasks shared.Tasks) {
	builder := &strings.Builder{}
	dumped, _ := tasks.DumpGroups()

	for _, g := range dumped {
		builder.WriteString("[\n")

		for _, t := range g {
			builder.WriteString("\t")
			builder.WriteString(t.Name)
			builder.WriteString("\n")
		}

		builder.WriteString("]\n")
	}

	glog.V(6).Infof("Executing Tasks:\n%s", builder.String())
}
