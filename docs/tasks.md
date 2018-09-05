# Tasks

Tasks for the installer should be small single steps (which do NOT correlate with the wizard steps).
For each task a correspondig go file should be placed under /pkg/tasks/.
Each task consists of a struct defining static parameters we know ahead of execution (e.g. for which node), the struct should implement the Task-interface and embed the BaseTask struct.

Example:
```go
type SomeAwesomeTask struct {
	shared.BaseTask
	nodeName string
}

func (t *SomeAwesomeTask) Execute(ctx *shared.Context) error {
	var node Node

	for _, n := range ctx.Nodes {
		if n.Name == t.nodeName {
			return nil
		}
	}

	return errors.New("no node found")
}
```

The passed Context contains data which can be shared between nodes (read-write).
Here you should remember, that you can only access data from previous tasks, not from tasks running at the same time.
The execution flow of tasks (as in: which tasks will be run before my task, which at the same time, etc.) correlates to the dependencies of the task.

When you want to add a new task, you'll have to implement the interface and extend the `setupTasks`-function of `/pkg/command/install.go`.

E.g. our SomeAwesomeTask depends on having a node provisioned:
```go
	provisionNode := addTask(&tasks.NodeProvisionTask{nodeName: "foobar"})

	awesomeTasks := shared.Tasks{
		addTask(&tasks.SomeAwesomeTask{nodeName: "foo"}, provisionNode),
		addTask(&tasks.SomeAwesomeTask{nodeName: "bar"}, provisionNode),
	}

	doSomethingElse := addTask(&tasks.SomeOtherTask{}, awesomeTasks...)
```

In this case our execution flow will look as following:
```

               provisionNode
                     |
           +---------+---------+
           |                   |
 doSomethingAwesome1   doSomethingAwesome2
           |                   |
           +---------+---------+
                     |
              doSomethingElse   

```

For debugging purposes the actual execution flow will be printed before executing when verbosity level 6 is defined.