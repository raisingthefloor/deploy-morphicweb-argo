package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func print_usage() {
	fmt.Println(`Usage:
	go run update-env-var.go [filepath] [scaffold key] [new value]

Update a kubernetes manifest using "scaffolding" comments, which allows us to preserve existing comments.

A scaffolding comment looks like...
          - name: MYVAL
            value: value  # scaffold:update-env-var:MYVAL

Example:
	go run update-env-var.go environments/dev/patches/set-env.yaml MYVAL my-new-value`)
}

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		print_usage()
		fmt.Println("must provide exactly 3 arguments")
		os.Exit(1)
	}

	err := do(args[0], args[1], args[2])
	if err != nil {
		fmt.Println("failed to update file with new value: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("success")
}

// Find the scaffold line, replace the contents, write the file back with the change.
func do(targetFilepath, scaffoldKey, newValue string) error {
	data, err := ioutil.ReadFile(targetFilepath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")

	fullScaffoldKey := "scaffold:update-env-var:" + scaffoldKey

	found := false
	for i, line := range lines {
		if !strings.Contains(line, fullScaffoldKey) {
			continue
		}

		found = true
		lines[i] = fmt.Sprintf(`            value: %s  # %s`, newValue, fullScaffoldKey)
	}

	if !found {
		return fmt.Errorf("failed to find scaffold key in file: %s", fullScaffoldKey)
	}

	newData := strings.Join(lines, "\n")

	return ioutil.WriteFile(targetFilepath, []byte(newData), 0611)
}
