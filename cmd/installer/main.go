package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"github.com/kubermatic/kubermatic-installer/pkg/command"
)

func main() {
	var wizard bool
	var install bool
	var manifestFile string

	flag.BoolVar(&wizard, "wizard", false, "Specify when you want to start the wizard.")
	flag.BoolVar(&install, "install", false, "Specify when you want to start the installer.")
	flag.StringVar(&manifestFile, "manifest", "", "file path to the manifest")

	flag.Parse()

	if wizard && install {
		glog.Error("only specify -wizard OR -install")
		os.Exit(1)
	}

	if wizard {
		err := command.WizardCommand()
		if err != nil {
			glog.Errorf("Error in wizard: %s", err)
			os.Exit(1)
		}
	} else if install {
		if manifestFile == "" {
			glog.Error("Please specify -manifest")
			os.Exit(1)
		}

		manifestContent, err := ioutil.ReadFile(manifestFile)
		if err != nil {
			glog.Errorf("Couldn't read manifest file: %v", err)
			os.Exit(1)
		}

		err = command.InstallCommand(manifestContent)
		if err != nil {
			glog.Errorf("Error in installer: %s", err)
			os.Exit(2)
		}
	} else {
		glog.Error("no command specified. Use -wizard or -install")
		os.Exit(1)
	}
}
