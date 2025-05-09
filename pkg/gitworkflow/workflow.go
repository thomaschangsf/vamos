package gitworkflow

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// WorkflowManager handles git workflow operations
type WorkflowManager struct {
	// Add configuration fields as needed
}

// NewWorkflowManager creates a new WorkflowManager instance
func NewWorkflowManager() *WorkflowManager {
	return &WorkflowManager{}
}

// SyncWithRemote syncs the current branch with remote
func (wm *WorkflowManager) SyncWithRemote() error {
	// First fetch to get latest changes without merging
	if err := wm.fetchOrigin(); err != nil {
		return fmt.Errorf("failed to fetch from origin: %w", err)
	}

	// Get current branch name
	currentBranch, err := wm.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Check if there are uncommitted changes (ignoring untracked files)
	cmd := exec.Command("git", "diff-files", "--quiet")
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return fmt.Errorf("uncommitted changes detected. Please commit or stash your changes before syncing")
		}
		return fmt.Errorf("failed to check git status: %w", err)
	}

	// Check if there are staged changes
	cmd = exec.Command("git", "diff-index", "--quiet", "--cached", "HEAD")
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return fmt.Errorf("staged changes detected. Please commit your changes before syncing")
		}
		return fmt.Errorf("failed to check git status: %w", err)
	}

	// Check if local branch has diverged from remote
	cmd = exec.Command("git", "rev-list", "--left-right", "--count", fmt.Sprintf("origin/%s...%s", currentBranch, currentBranch))
	output, err := cmd.Output()
	if err != nil {
		// If the remote branch doesn't exist yet, that's okay - we'll create it later
		if strings.Contains(err.Error(), "unknown revision") {
			return nil
		}
		return fmt.Errorf("failed to check branch divergence: %w", err)
	}

	// Parse the output (format: "X\tY" where X is commits ahead, Y is commits behind)
	parts := strings.Fields(string(output))
	if len(parts) != 2 {
		return fmt.Errorf("unexpected output from rev-list command")
	}

	ahead, _ := strconv.Atoi(parts[0])
	behind, _ := strconv.Atoi(parts[1])

	if ahead > 0 {
		return fmt.Errorf("local branch is ahead of remote by %d commits. Please push your changes first using 'git push' or 'make story-push'", ahead)
	}

	if behind > 0 {
		// Pull changes from remote
		cmd = exec.Command("git", "pull", "--rebase")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to pull changes: %w", err)
		}
	}

	return nil
}

// getDefaultBranch determines whether the repository uses 'main' or 'master'
func (wm *WorkflowManager) getDefaultBranch() (string, error) {
	// Try to get the ref for main
	cmdMain := exec.Command("git", "rev-parse", "--verify", "refs/heads/main")
	if err := cmdMain.Run(); err == nil {
		return "main", nil
	}

	// Try to get the ref for master
	cmdMaster := exec.Command("git", "rev-parse", "--verify", "refs/heads/master")
	if err := cmdMaster.Run(); err == nil {
		return "master", nil
	}

	return "", fmt.Errorf("neither 'main' nor 'master' branch found")
}

// SyncMainBranch syncs the main branch with remote (works with either main or master)
func (wm *WorkflowManager) SyncMainBranch() error {
	// Determine which default branch is used
	defaultBranch, err := wm.getDefaultBranch()
	if err != nil {
		return err
	}

	// Checkout main/master branch
	if err := wm.checkoutBranch(defaultBranch); err != nil {
		return fmt.Errorf("failed to checkout %s branch: %w", defaultBranch, err)
	}

	// Pull latest changes
	if err := wm.pullLatest(); err != nil {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}

	return nil
}

// ResolveConflictsRebase resolves conflicts by rebasing onto origin/main
func (wm *WorkflowManager) ResolveConflictsRebase() error {
	// Fetch latest changes
	if err := wm.fetchOrigin(); err != nil {
		return fmt.Errorf("failed to fetch from origin: %w", err)
	}

	// Rebase onto origin/main
	cmd := exec.Command("git", "rebase", "origin/main")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to rebase onto origin/main: %w", err)
	}

	return nil
}

// ResolveConflictsMerge resolves conflicts by merging origin/main
func (wm *WorkflowManager) ResolveConflictsMerge() error {
	// Fetch latest changes
	if err := wm.fetchOrigin(); err != nil {
		return fmt.Errorf("failed to fetch from origin: %w", err)
	}

	// Merge origin/main
	cmd := exec.Command("git", "merge", "origin/main")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to merge origin/main: %w", err)
	}

	return nil
}

// Helper function to fetch from origin
func (wm *WorkflowManager) fetchOrigin() error {
	cmd := exec.Command("git", "fetch", "origin")
	return cmd.Run()
}

// Helper function to checkout a branch
func (wm *WorkflowManager) checkoutBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", branchName)
	return cmd.Run()
}

// Helper function to pull latest changes
func (wm *WorkflowManager) pullLatest() error {
	cmd := exec.Command("git", "pull")
	return cmd.Run()
}

// CreateStoryBranch creates a new story branch from the main branch
func (wm *WorkflowManager) CreateStoryBranch(storyID string, description string) error {
	// Format the branch name
	var branchName string
	if description != "" {
		branchName = fmt.Sprintf("W-%s-%s", storyID, strings.ToLower(strings.ReplaceAll(description, " ", "-")))
	} else {
		branchName = fmt.Sprintf("W-%s", storyID)
	}

	// Ensure we're on main branch
	defaultBranch, err := wm.getDefaultBranch()
	if err != nil {
		return err
	}

	if err := wm.checkoutBranch(defaultBranch); err != nil {
		return fmt.Errorf("failed to checkout %s branch: %w", defaultBranch, err)
	}

	// Pull latest changes
	if err := wm.pullLatest(); err != nil {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}

	// Create and checkout new story branch
	cmd := exec.Command("git", "checkout", "-b", branchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create story branch: %w", err)
	}

	return nil
}

// CommitChanges creates a commit with a formatted message
func (wm *WorkflowManager) CommitChanges(scope string, description string) error {
	// Format the commit message
	commitMessage := fmt.Sprintf("feat(%s): %s", scope, description)

	// Add all changes
	addCmd := exec.Command("git", "add", ".")
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}

	// Create the commit
	commitCmd := exec.Command("git", "commit", "-m", commitMessage)
	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	return nil
}

// PushStoryBranch pushes the current story branch to remote
func (wm *WorkflowManager) PushStoryBranch() error {
	// Get current branch name
	branchName, err := wm.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Push the branch to remote
	cmd := exec.Command("git", "push", "-u", "origin", branchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push branch: %w", err)
	}

	return nil
}

// validateBranchName checks if the branch name follows the convention
func (wm *WorkflowManager) validateBranchName(branchName string) error {
	if !strings.HasPrefix(branchName, "W-") {
		return fmt.Errorf("branch name must follow the format: W-STORY_ID (e.g., W-123)")
	}

	// Check if there's a story ID after the "W-" prefix
	parts := strings.SplitN(branchName[2:], "-", 2) // Skip "W-" prefix and split the rest
	if len(parts) == 0 || parts[0] == "" {
		return fmt.Errorf("branch name must include a story ID")
	}

	return nil
}

// CreateFeatureBranch creates a new feature branch from the main branch
func (wm *WorkflowManager) CreateFeatureBranch(branchName string) error {
	// Validate branch name format
	if err := wm.validateBranchName(branchName); err != nil {
		return err
	}

	// Ensure we're on main branch
	if err := wm.checkoutBranch("main"); err != nil {
		return fmt.Errorf("failed to checkout main branch: %w", err)
	}

	// Pull latest changes
	if err := wm.pullLatest(); err != nil {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}

	// Create and checkout new branch
	cmd := exec.Command("git", "checkout", "-b", branchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

// CreatePullRequest prepares a branch for pull request submission
func (wm *WorkflowManager) CreatePullRequest(branchName string) error {
	// Ensure we're on the feature branch
	if err := wm.checkoutBranch(branchName); err != nil {
		return fmt.Errorf("failed to checkout feature branch: %w", err)
	}

	// Push the branch to remote
	cmd := exec.Command("git", "push", "-u", "origin", branchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push branch: %w", err)
	}

	return nil
}

// GetCurrentBranch returns the name of the current git branch
func (wm *WorkflowManager) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// UndoLastCommit undoes the last commit, keeping changes in working directory
func (wm *WorkflowManager) UndoLastCommit() error {
	cmd := exec.Command("git", "reset", "HEAD~1")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to undo last commit: %w", err)
	}
	return nil
}

// UndoLastCommitHard undoes the last commit and discards changes
func (wm *WorkflowManager) UndoLastCommitHard() error {
	cmd := exec.Command("git", "reset", "--hard", "HEAD~1")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to undo last commit (hard): %w", err)
	}
	return nil
}

// RevertCommit creates a new commit that undoes the changes of a specific commit
func (wm *WorkflowManager) RevertCommit(commitHash string) error {
	cmd := exec.Command("git", "revert", commitHash)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to revert commit %s: %w", commitHash, err)
	}
	return nil
}

// CreateTag creates a new tag at the current commit
func (wm *WorkflowManager) CreateTag(version, message string) error {
	cmd := exec.Command("git", "tag", "-a", version, "-m", message)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tag %s: %w", version, err)
	}
	return nil
}

// PushTag pushes a specific tag to remote
func (wm *WorkflowManager) PushTag(version string) error {
	cmd := exec.Command("git", "push", "origin", version)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push tag %s: %w", version, err)
	}
	return nil
}

// PushAllTags pushes all tags to remote
func (wm *WorkflowManager) PushAllTags() error {
	cmd := exec.Command("git", "push", "origin", "--tags")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push all tags: %w", err)
	}
	return nil
}

// GetLastCommitHash returns the hash of the last commit
func (wm *WorkflowManager) GetLastCommitHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get last commit hash: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// PrintExample prints out an end-to-end example of the git workflow
func (wm *WorkflowManager) PrintExample() {
	fmt.Print(`
End-to-End Git Workflow Example
==============================

1. Start a new story (with description):
   vamosGitWF story-start --id "123" --description "add-user-authentication"

   Or without description:
   vamosGitWF story-start --id "123"

2. Make some changes and commit them:
   vamosGitWF story-commit --scope "auth" --description "implement basic login flow"

3. If you need to undo your last commit (keeping changes):
   vamosGitWF undo

4. If you need to undo your last commit (discarding changes):
   vamosGitWF undo --hard

5. If you need to revert a specific commit:
   vamosGitWF revert --commit "abc123"

6. When you're ready to sync with remote changes:
   vamosGitWF sync

7. If there are conflicts, resolve them:
   vamosGitWF resolve

8. When the story is complete, create a version tag:
   vamosGitWF tag --version "v1.0.0" --message "Initial release" --push

9. If you need to revert to a previous tag:
    git checkout v1.0.0  # Checkout the tag
    vamosGitWF story-start --id "124"  # Create a new branch from the tag

10. Push your changes to remote:
    vamosGitWF story-push

11. Finally, sync the main branch:
    vamosGitWF sync --main

Best Practices:
- Always start from the main branch
- Commit frequently with clear, descriptive messages
- Push your changes regularly
- Use sync before starting new work
- Resolve conflicts as soon as they appear
- Tag releases when features are complete
- Use tags as stable points to revert to if needed
`)
}
