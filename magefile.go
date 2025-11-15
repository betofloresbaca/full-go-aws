//go:build mage
// +build mage

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
)

var LambdasDir = "cmd/lambdas"
var BinDir = "bin"

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

// Clean up after yourself
func Clean() error {
	fmt.Println("Cleaning...")
	fmt.Println("Deleting bin")
	if err := os.RemoveAll("bin"); err != nil {
		return err
	}
	fmt.Println("Deleting cdk.out")
	if err := os.RemoveAll("cdk.out"); err != nil {
		return err
	}
	return nil
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	fmt.Println("Building...")
	fmt.Println("Building Lambdas...")
	return buildLambdas()
}

func Deploy() error {
	mg.Deps(Build)
	fmt.Println("Deploying CDK Stack...")
	cdkDeployCmd := exec.Command("cdk", "deploy", "--all", "--require-approval", "never")
	cdkDeployCmd.Stdout = os.Stdout
	cdkDeployCmd.Stderr = os.Stderr
	return cdkDeployCmd.Run()
}

func Bootstrap() error {
	mg.Deps(Build)
	fmt.Println("Bootstrapping CDK...")
	cdkBootstrapCmd := exec.Command("cdk", "bootstrap")
	cdkBootstrapCmd.Stdout = os.Stdout
	cdkBootstrapCmd.Stderr = os.Stderr
	return cdkBootstrapCmd.Run()
}

func Destroy() error {
	fmt.Println("Destroying CDK Stack...")
	cdkDestroyCmd := exec.Command("cdk", "destroy", "--all", "--force")
	cdkDestroyCmd.Stdout = os.Stdout
	cdkDestroyCmd.Stderr = os.Stderr
	return cdkDestroyCmd.Run()
}

func buildLambdas() error {
	lambdas := getLambdaNames()
	for _, name := range lambdas {
		fmt.Println("Building lambda:", name)
		goBuildCmd := exec.Command("go", "build", "-o", "../../../"+BinDir+"/"+name, ".")
		goBuildCmd.Dir = LambdasDir + "/" + name
		goBuildCmd.Stdout = os.Stdout
		goBuildCmd.Stderr = os.Stderr
		goBuildCmd.Env = append(os.Environ(),
			"GOOS=linux",
			"GOARCH=arm64",
		)
		if err := goBuildCmd.Run(); err != nil {
			return fmt.Errorf("Failed to build %s: %w", name, err)
		}
		zipPath := BinDir + "/" + name + ".zip"
		binaryPath := BinDir + "/" + name
		if err := createZip(binaryPath, zipPath); err != nil {
			return fmt.Errorf("Failed to create zip for %s: %w", name, err)
		}
		fmt.Println("Created zip:", zipPath)
		os.Remove(binaryPath)
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

func createZip(binaryPath, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	binary, err := os.Open(binaryPath)
	if err != nil {
		return err
	}
	defer binary.Close()

	// IMPORTANT: For PROVIDED_AL2023 runtime, the file should be named 'bootstrap'.
	// For PROVIDED_AL2 runtime, it should be named 'main'.
	w, err := archive.Create("bootstrap")
	if err != nil {
		return err
	}

	_, err = io.Copy(w, binary)
	return err
}
