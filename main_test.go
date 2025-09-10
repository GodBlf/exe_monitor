package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestMa(t *testing.T) {
	for {
		output, _ := exec.Command("tasklist").Output()
		if strings.Contains(string(output), "AweSun.exe") {
			fmt.Println("exist")
		} else {
			fmt.Println("not exist")
		}
	}

}
