package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	flagsRE = regexp.MustCompile(`\sflags=0x[0-9a-bA-B]+\((.*)\)\s`)
)

func main() {
	apps, err := os.ReadDir("/Applications")
	if err != nil {
		panic(err)
	}

	for _, app := range apps {
		if strings.Contains(app.Name(), ".app") {
			appPath := filepath.Join("/Applications", app.Name())
			buf := new(bytes.Buffer)
			cmd := exec.Command("codesign", "-dv", appPath)
			cmd.Stderr = buf
			if err := cmd.Run(); err != nil {
				continue
			}
			flags := getFlags(buf.String())
			found := false
			for _, flag := range flags {
				if flag == "runtime" {
					found = true
				}
			}
			if !found {
				fmt.Printf("[*] App %s is missing \"runtime\" flag; flags=%v\n", app.Name(), flags)
			}
		}
	}
}

func getFlags(codesignOutput string) []string {
	matches := flagsRE.FindStringSubmatch(codesignOutput)
	if len(matches) > 1 {
		flagsString := strings.Split(matches[1], ",")
		return flagsString
	} else {
		return []string{}
	}
}
