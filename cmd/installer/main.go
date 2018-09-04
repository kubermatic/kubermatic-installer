package main

import (
	"flag"
	"fmt"
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
		glog.Fatalf("only specify -wizard OR -install")
	}

	if wizard {
		err := command.WizardCommand()
		if err != nil {
			fmt.Printf("Error in wizard: %s\n", err)
			os.Exit(1)
		}
	} else if install {
		if manifestFile == "" {
			fmt.Println("Please specify -manifest")
			os.Exit(1)
		}

		manifestContent, err := ioutil.ReadFile(manifestFile)
		if err != nil {
			fmt.Printf("Couldn't read manifest file: %v\n", err)
			os.Exit(1)
		}

		err = command.InstallCommand(manifestContent)
		if err != nil {
			fmt.Printf("Error in installer: %s\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("no command specified. Use -wizard or -install\n")
		os.Exit(1)
	}
}
