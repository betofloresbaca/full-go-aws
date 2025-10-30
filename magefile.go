//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	// "github.com/magefile/mage/mg"
)

var LambdasDir = "cmd/lambdas"
var BinDir = "bin"

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	fmt.Println("Building...")
	fmt.Println("Building Lambdas...")
	return buildLambdas()
}

// Clean up after yourself
func Clean() error {
	fmt.Println("Cleaning...")
	fmt.Println("Deleting bin")
	return os.RemoveAll("bin")
}

func buildLambdas() error {
	lambdas := getLambdaNames()
	for _, name := range lambdas {
		fmt.Println("Building lambda:", name)
		cmd := exec.Command("go", "build", "-o", "../../../"+BinDir+"/"+name, ".")
		cmd.Dir = LambdasDir + "/" + name
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(os.Environ(),
			"GOOS=linux",
			"GOARCH=amd64",
		)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Failed to build %s: %w", name, err)
		}
	}
	return nil
}

func getLambdaNames() []string {
	lambdas := make([]string, 0)
	entries, _ := os.ReadDir(LambdasDir)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		lambdas = append(lambdas, e.Name())
	}
	return lambdas
}
