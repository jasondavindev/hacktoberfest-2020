package listener

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jasondavindev/hacktoberfest-2020/command"
	"gopkg.in/fsnotify.v1"
)

type ChangesListener struct {
	watcher             *fsnotify.Watcher
	excludedDirectories []string
	jobRunner           command.JobRunner
}

func CreateChangesListener(excludedDirectories string, jobRunner command.JobRunner) ChangesListener {
	listener := ChangesListener{}
	listener.watcher = CreateWatcher()
	listener.excludedDirectories = splitExcludedFiles(excludedDirectories)
	listener.jobRunner = jobRunner

	return listener
}

func CreateWatcher() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	return watcher
}

func isModifiedFile(e fsnotify.Event) bool {
	return e.Op == fsnotify.Create || e.Op == fsnotify.Remove || e.Op == fsnotify.Write
}

func splitExcludedFiles(excludedFiles string) []string {
	return strings.Split(excludedFiles, ",")
}

func (cl *ChangesListener) CloseWatcher() {
	cl.watcher.Close()
}

func (cl *ChangesListener) ListenEvents() {
	for {
		select {
		case event, ok := <-cl.watcher.Events:
			if !ok {
				return
			}

			if cl.isExcludedFile(event.Name) || isHiddenFile(event.Name) {
				continue
			}

			cl.EventHandler(event)
		case err, ok := <-cl.watcher.Errors:
			if !ok {
				return
			}

			log.Println("error:", err)
		}
	}
}

func (cl *ChangesListener) isExcludedFile(absoluteFile string) bool {
	fileName := filepath.Base(absoluteFile)

	for _, file := range cl.excludedDirectories {
		if fileName == file {
			return true
		}
	}

	return false
}

func (cl *ChangesListener) SetupDirectoriesToWatch(directory string) {
	directories, err := findSubDirectories(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range directories {
		err = cl.watcher.Add(file)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (cl *ChangesListener) EventHandler(event fsnotify.Event) bool {
	if isModifiedFile(event) {
		formatResponse(cl.jobRunner.RunJobs())
		return true
	}

	return false
}

func isHiddenFile(path string) bool {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		log.Panic(err)
	}

	fileName := filepath.Base(absolutePath)

	return filepath.HasPrefix(absolutePath, ".") || fileName[0:1] == "."
}

func findSubDirectories(directory string) ([]string, error) {
	paths := []string{}

	return paths, filepath.Walk(directory, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if isHiddenFile(newPath) {
				return filepath.SkipDir
			}
			paths = append(paths, newPath)
		}

		return nil
	})
}

func formatResponse(res []string) {
	for _, str := range res {
		fmt.Println(str)
	}
}

func RunCommandsAndFormatResponse(jobRunner *command.JobRunner) {
	formatResponse(jobRunner.RunJobs())
}
