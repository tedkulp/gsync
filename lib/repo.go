package lib

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

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
		//fmt.Fprintln(os.Stderr, "There was an error running git remove command: ", err)
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

func GetHasRemote(repoLocation string, remoteName string) bool {
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
		"remote",
		"show",
		remoteName,
	}

	if _, stderr, _ = RunCommand(cmdName, cmdArgs...); stderr != "" {
		fmt.Fprintln(os.Stderr, "There was an error running git remote command: ", stderr)
	}

	return exitCode == 0
}

func GitPush(repoLocation string, remoteName string) {
	var (
		err error
	)

	workTree := filepath.Join(repoLocation)
	gitDir := filepath.Join(repoLocation, ".git")

	cmdName := "git"
	cmdArgs := []string{
		"--work-tree=" + workTree,
		"--git-dir=" + gitDir,
		"push",
		remoteName,
		"master",
	}

	if _, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git push command: ", err)
	}
}

func GitPull(repoLocation string, remoteName string) {
	var (
		err error
	)

	workTree := filepath.Join(repoLocation)
	gitDir := filepath.Join(repoLocation, ".git")

	cmdName := "git"
	cmdArgs := []string{
		"--work-tree=" + workTree,
		"--git-dir=" + gitDir,
		"pull",
		remoteName,
		"master",
	}

	if _, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git pull command: ", err)
	}
}
