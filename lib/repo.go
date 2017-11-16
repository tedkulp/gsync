package lib

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const defaultFailedCode = 1

func CopyFileToRepo(filename string, hostname string, repoLocation string) {
	directoryPath := filepath.Join(repoLocation, hostname, filepath.Dir(filename))
	pathErr := os.MkdirAll(directoryPath, 0755)
	if pathErr != nil {
		fmt.Println(pathErr)
		return
	}

	CopyFile(filename, filepath.Join(directoryPath, filepath.Base(filename)))
}

func CopyFile(source string, dest string) {
	fmt.Println("copying: " + source + " => " + dest)

	from, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}

func GitAdd(filename string, hostname string, repoLocation string) {
	var (
		err error
	)

	fileToAdd := filepath.Join(repoLocation, hostname, filename)
	workTree := filepath.Join(repoLocation)
	gitDir := filepath.Join(repoLocation, ".git")

	cmdName := "git"
	cmdArgs := []string{
		"--work-tree=" + workTree,
		"--git-dir=" + gitDir,
		"add",
		fileToAdd,
	}

	if _, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git add command: ", err)
	}
}

func GitRemove(filename string, hostname string, repoLocation string) {
	var (
		err error
	)

	fileToRemove := filepath.Join(repoLocation, hostname, filename)
	workTree := filepath.Join(repoLocation)
	gitDir := filepath.Join(repoLocation, ".git")

	cmdName := "git"
	cmdArgs := []string{
		"--work-tree=" + workTree,
		"--git-dir=" + gitDir,
		"rm",
		fileToRemove,
	}

	fmt.Println(cmdArgs)

	if _, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git remove command: ", err)
	}
}

func GitCommit(hostname string, repoLocation string) {
	var (
		err error
	)

	workTree := filepath.Join(repoLocation)
	gitDir := filepath.Join(repoLocation, ".git")

	cmdName := "git"
	cmdArgs := []string{
		"--work-tree=" + workTree,
		"--git-dir=" + gitDir,
		"commit",
		"-m",
		"Update files on host '" + hostname + "'",
	}

	if _, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git commit command: ", err)
	}
}

func GitHasChangesToCommit(repoLocation string) bool {
	var (
		stderr   string
		exitCode int
	)

	workTree := filepath.Join(repoLocation)
	gitDir := filepath.Join(repoLocation, ".git")

	cmdName := "git"
	cmdArgs := []string{
		"--work-tree=" + workTree,
		"--git-dir=" + gitDir,
		"diff",
		"--cached",
		"--exit-code",
	}

	if _, stderr, exitCode = RunCommand(cmdName, cmdArgs...); stderr != "" {
		fmt.Fprintln(os.Stderr, "There was an error running git commit command: ", stderr)
	}

	return exitCode > 0
}

func RunCommand(name string, args ...string) (stdout string, stderr string, exitCode int) {
	log.Println("run command:", name, args)
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			log.Printf("Could not get exit code for failed program: %v, %v", name, args)
			exitCode = defaultFailedCode
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	log.Printf("command result, stdout: %v, stderr: %v, exitCode: %v", stdout, stderr, exitCode)
	return
}
