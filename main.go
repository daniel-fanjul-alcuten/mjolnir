package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

const GOPATH = "GOPATH"
const GOPATH_PREFIX = GOPATH + "="

// the configuration file
const JSON_FILE = "mjolnir.json"

// the configuration extracted from the JSON_FILE
type Config struct {
	GoDir           string
	GoPaths         []string
	TestPackages    []string
	InstallPackages []string
	Exec            string
}

func main() {
	log.SetFlags(0)

	file, err := os.Open(JSON_FILE)
	if err != nil {
		log.Fatalln("error:", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("error:", err)
	}

	goDir, err := filepath.Abs(config.GoDir)
	if err != nil {
		log.Fatalln("error:", err)
	}

	gopaths := make([]string, len(config.GoPaths))
	for i, path := range config.GoPaths {
		gopaths[i], err = filepath.Abs(path)
		if err != nil {
			log.Fatalln("error:", err)
		}
	}

	environ := os.Environ()
	gopath := strings.Join(gopaths, string(os.PathListSeparator))
	ok := false
	for index, value := range environ {
		if strings.HasPrefix(value, GOPATH_PREFIX) {
			if len(gopath) > 0 {
				environ[index] = GOPATH_PREFIX + strings.Join([]string{gopath, value[len(GOPATH_PREFIX):]}, string(os.PathListSeparator))
			}
			ok = true
			break
		}
	}
	if !ok && len(gopath) > 0 {
		environ = append(environ, GOPATH_PREFIX+gopath)
	}

	if len(config.TestPackages) > 0 {
		args := append([]string{"test"}, config.TestPackages...)
		line := []interface{}{"+", "go"}
		for _, arg := range args {
			line = append(line, arg)
		}
		log.Println(line...)
		test := exec.Command("go", args...)
		test.Dir = goDir
		test.Env = environ
		test.Stdout = os.Stdout
		test.Stderr = os.Stderr
		err = test.Run()
		if err != nil {
			log_err(err)
		}
	}

	if len(config.InstallPackages) > 0 {
		args := append([]string{"install"}, config.InstallPackages...)
		line := []interface{}{"+", "go"}
		for _, arg := range args {
			line = append(line, arg)
		}
		log.Println(line...)
		build := exec.Command("go", args...)
		build.Dir = goDir
		build.Env = environ
		build.Stdout = os.Stdout
		build.Stderr = os.Stderr
		err = build.Run()
		if err != nil {
			log_err(err)
		}
	}

	if config.Exec != "" {
		execc, err := filepath.Abs(config.Exec)
		if err != nil {
			log.Fatalln("error:", err)
		}
		args := os.Args[1:]
		line := []interface{}{"+", execc}
		for _, arg := range args {
			line = append(line, arg)
		}
		log.Println(line...)
		run := exec.Command(execc, args...)
		run.Stdout = os.Stdout
		run.Stderr = os.Stderr
		err = run.Run()
		if err != nil {
			log_err(err)
		}
	}
}

func log_err(err error) {
	if exiterr, ok := err.(*exec.ExitError); ok {
		sys := exiterr.ProcessState.Sys()
		if status, ok := sys.(syscall.WaitStatus); ok {
			os.Exit(status.ExitStatus())
		}
	}
	log.Fatalln("error:", err)
}
