package kirin

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func CreateProject(appName, moduleName, frontendChoice string) (err error) {
	cloneUrl := "https://github.com/thetnaingtn/boilerplate"

	wd, _ := os.Getwd()
	projectPath := fmt.Sprintf("%s%c%s", wd, os.PathSeparator, appName)

	defer func() {
		if err != nil {
			os.RemoveAll(projectPath)
		}
	}()

	var git string
	git, err = exec.LookPath("git")

	if err != nil {
		return fmt.Errorf("git is not installed or not found in PATH")
	}

	cmd := exec.Command(git, "clone", "-b", fmt.Sprintf("frontend/%s", frontendChoice), cloneUrl, projectPath)

	if err = cmd.Run(); err != nil {
		return
	}

	if err = replace(projectPath, "go.mod", "boilerplate", moduleName); err != nil {
		return
	}

	if err = replace(projectPath, "*.go", "boilerplate", moduleName); err != nil {
		return
	}

	return nil
}

func replace(path, pattern, old, new string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return replaceWalkFn(path, info, pattern, []byte(old), []byte(new))
	})
}

func replaceWalkFn(path string, info os.FileInfo, pattern string, old, new []byte) (err error) {
	var matched bool
	if matched, err = filepath.Match(pattern, info.Name()); err != nil {
		return
	}

	if matched {
		cleanedPath := filepath.Clean(path)

		var oldContent []byte
		if oldContent, err = os.ReadFile(cleanedPath); err != nil {
			return
		}

		if err = os.WriteFile(cleanedPath, bytes.ReplaceAll(oldContent, old, new), 0); err != nil {
			return
		}
	}

	return
}
