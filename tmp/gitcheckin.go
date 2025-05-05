// gitcheckin.go
// A simple CLI tool to streamline Git check-ins: stage changes, compose commit message, and push.

package main

import (
    "bufio"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "strings"
)

func main() {
    var all bool
    var msgFlag string
    var noPush bool

    flag.BoolVar(&all, "a", false, "Stage all modified/deleted/new files (git add --all)")
    flag.StringVar(&msgFlag, "m", "", "Commit message. If empty, opens $EDITOR")
    flag.BoolVar(&noPush, "no-push", false, "Do not push after commit")
    flag.Parse()

    if all {
        runGit("add", "--all")
    }

    var commitMsg string
    if msgFlag != "" {
        commitMsg = msgFlag
    } else {
        commitMsg = getMessageFromEditor()
        if strings.TrimSpace(commitMsg) == "" {
            fmt.Fprintln(os.Stderr, "Aborting commit due to empty commit message.")
            os.Exit(1)
        }
    }

    runGit("commit", "-m", commitMsg)

    if !noPush {
        runGit("push")
    }
}

// runGit executes git command with args and streams output.
func runGit(args ...string) {
    cmd := exec.Command("git", args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    if err := cmd.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "git %s failed: %v\n", strings.Join(args, " "), err)
        os.Exit(1)
    }
}

// getMessageFromEditor opens user's EDITOR for commit message.
func getMessageFromEditor() string {
    editor := os.Getenv("EDITOR")
    if editor == "" {
        editor = "vi"
    }

    tmpFile, err := ioutil.TempFile("", "GITCHECKIN_*.tmp")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create temp file: %v\n", err)
        os.Exit(1)
    }
    tmpName := tmpFile.Name()
    tmpFile.Close()
    defer os.Remove(tmpName)

    cmd := exec.Command(editor, tmpName)
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "%s failed: %v\n", editor, err)
        os.Exit(1)
    }

    data, err := ioutil.ReadFile(tmpName)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to read temp file: %v\n", err)
        os.Exit(1)
    }
    // Strip out comment lines
    var lines []string
    scanner := bufio.NewScanner(strings.NewReader(string(data)))
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "#") {
            continue
        }
        lines = append(lines, line)
    }
    return strings.Join(lines, "\n")
}