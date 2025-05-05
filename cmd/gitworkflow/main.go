package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/thomaschangsf/vamos/pkg/gitworkflow"
)

func main() {
	// Create workflow manager
	wm := gitworkflow.NewWorkflowManager()

	// Define subcommands
	storyStartCmd := flag.NewFlagSet("story-start", flag.ExitOnError)
	storyID := storyStartCmd.String("id", "", "Story ID (required)")
	description := storyStartCmd.String("description", "", "Story description (required)")

	storyCommitCmd := flag.NewFlagSet("story-commit", flag.ExitOnError)
	scope := storyCommitCmd.String("scope", "", "Commit scope (required)")
	commitDesc := storyCommitCmd.String("description", "", "Commit description (required)")

	undoCmd := flag.NewFlagSet("undo", flag.ExitOnError)
	hard := undoCmd.Bool("hard", false, "Discard changes (default: keep changes)")

	revertCmd := flag.NewFlagSet("revert", flag.ExitOnError)
	commitHash := revertCmd.String("commit", "", "Commit hash to revert (required)")

	tagCmd := flag.NewFlagSet("tag", flag.ExitOnError)
	version := tagCmd.String("version", "", "Version number (e.g., v1.0.3) (required)")
	tagMessage := tagCmd.String("message", "", "Tag message (required)")
	pushTag := tagCmd.Bool("push", false, "Push tag to remote")

	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)
	syncMain := syncCmd.Bool("main", false, "Sync main branch (default: sync current branch)")

	resolveCmd := flag.NewFlagSet("resolve", flag.ExitOnError)
	useRebase := resolveCmd.Bool("rebase", true, "Use rebase to resolve conflicts (default: true)")

	// Check if a subcommand was provided
	if len(os.Args) < 2 {
		fmt.Println("Expected one of: 'story-start', 'story-commit', 'story-push', 'undo', 'revert', 'tag', 'sync', 'resolve', or 'example'")
		os.Exit(1)
	}

	// Parse the subcommand
	switch os.Args[1] {
	case "example":
		wm.PrintExample()
		os.Exit(0)

	case "story-start":
		storyStartCmd.Parse(os.Args[2:])
		if *storyID == "" || *description == "" {
			fmt.Println("Error: --id and --description are required")
			storyStartCmd.PrintDefaults()
			os.Exit(1)
		}
		err := wm.CreateStoryBranch(*storyID, *description)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created and switched to branch: feat/story-%s-%s\n", *storyID, *description)

	case "story-commit":
		storyCommitCmd.Parse(os.Args[2:])
		if *scope == "" || *commitDesc == "" {
			fmt.Println("Error: --scope and --description are required")
			storyCommitCmd.PrintDefaults()
			os.Exit(1)
		}
		err := wm.CommitChanges(*scope, *commitDesc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Committed changes: feat(%s): %s\n", *scope, *commitDesc)

	case "story-push":
		err := wm.PushStoryBranch()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Pushed branch to remote")

	case "undo":
		undoCmd.Parse(os.Args[2:])
		var err error
		if *hard {
			err = wm.UndoLastCommitHard()
			fmt.Println("Undid last commit and discarded changes")
		} else {
			err = wm.UndoLastCommit()
			fmt.Println("Undid last commit, changes are in working directory")
		}
		if err != nil {
			log.Fatal(err)
		}

	case "revert":
		revertCmd.Parse(os.Args[2:])
		if *commitHash == "" {
			fmt.Println("Error: --commit is required")
			revertCmd.PrintDefaults()
			os.Exit(1)
		}
		err := wm.RevertCommit(*commitHash)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Reverted commit %s\n", *commitHash)

	case "tag":
		tagCmd.Parse(os.Args[2:])
		if *version == "" || *tagMessage == "" {
			fmt.Println("Error: --version and --message are required")
			tagCmd.PrintDefaults()
			os.Exit(1)
		}
		err := wm.CreateTag(*version, *tagMessage)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created tag %s: %s\n", *version, *tagMessage)

		if *pushTag {
			err = wm.PushTag(*version)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Pushed tag %s to remote\n", *version)
		}

	case "sync":
		syncCmd.Parse(os.Args[2:])
		var err error
		if *syncMain {
			err = wm.SyncMainBranch()
			fmt.Println("Synced main branch with remote")
		} else {
			err = wm.SyncWithRemote()
			fmt.Println("Synced current branch with remote")
		}
		if err != nil {
			log.Fatal(err)
		}

	case "resolve":
		resolveCmd.Parse(os.Args[2:])
		var err error
		if *useRebase {
			err = wm.ResolveConflictsRebase()
			fmt.Println("Resolved conflicts by rebasing onto origin/main")
		} else {
			err = wm.ResolveConflictsMerge()
			fmt.Println("Resolved conflicts by merging origin/main")
		}
		if err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Println("Expected one of: 'story-start', 'story-commit', 'story-push', 'undo', 'revert', 'tag', 'sync', 'resolve', or 'example'")
		os.Exit(1)
	}
}
