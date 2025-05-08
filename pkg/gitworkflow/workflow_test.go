package gitworkflow

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestNewWorkflowManager(t *testing.T) {
	wm := NewWorkflowManager()
	if wm == nil {
		t.Error("Expected WorkflowManager instance, got nil")
	}
}

func TestGetCurrentBranch(t *testing.T) {
	wm := NewWorkflowManager()
	branch, err := wm.GetCurrentBranch()
	if err != nil {
		t.Errorf("Unexpected error getting current branch: %v", err)
	}
	if branch == "" {
		t.Error("Expected non-empty branch name")
	}
}

func TestValidateBranchName(t *testing.T) {
	wm := NewWorkflowManager()

	tests := []struct {
		name    string
		branch  string
		wantErr bool
	}{
		{
			name:    "valid story branch",
			branch:  "W-123",
			wantErr: false,
		},
		{
			name:    "valid story branch with description",
			branch:  "W-123-add-login-button",
			wantErr: false,
		},
		{
			name:    "invalid prefix",
			branch:  "X-123",
			wantErr: true,
		},
		{
			name:    "missing story ID",
			branch:  "W-",
			wantErr: true,
		},
		{
			name:    "no prefix",
			branch:  "123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wm.validateBranchName(tt.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateBranchName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStoryBranchNameFormat(t *testing.T) {
	storyID := "456"
	description := "Chat UI"

	// Test with description
	branchNameWithDesc := fmt.Sprintf("W-%s-%s", storyID, strings.ToLower(strings.ReplaceAll(description, " ", "-")))
	expectedWithDesc := "W-456-chat-ui"
	if branchNameWithDesc != expectedWithDesc {
		t.Errorf("Expected branch name with description %s, got %s", expectedWithDesc, branchNameWithDesc)
	}

	// Test without description
	branchNameWithoutDesc := fmt.Sprintf("W-%s", storyID)
	expectedWithoutDesc := "W-456"
	if branchNameWithoutDesc != expectedWithoutDesc {
		t.Errorf("Expected branch name without description %s, got %s", expectedWithoutDesc, branchNameWithoutDesc)
	}
}

func TestGetDefaultBranch(t *testing.T) {
	wm := NewWorkflowManager()

	// Save current branch to restore later
	currentBranch, err := wm.GetCurrentBranch()
	if err != nil {
		t.Fatalf("Failed to get current branch: %v", err)
	}

	// Get the default branch
	defaultBranch, err := wm.getDefaultBranch()
	if err != nil {
		t.Fatalf("Failed to get default branch: %v", err)
	}

	// Verify that the default branch is either main or master
	if defaultBranch != "main" && defaultBranch != "master" {
		t.Errorf("Expected default branch to be 'main' or 'master', got %s", defaultBranch)
	}

	// Verify that the branch actually exists
	cmd := exec.Command("git", "rev-parse", "--verify", fmt.Sprintf("refs/heads/%s", defaultBranch))
	if err := cmd.Run(); err != nil {
		t.Errorf("Default branch %s does not exist", defaultBranch)
	}

	// Restore original branch if different
	if currentBranch != defaultBranch {
		if err := wm.checkoutBranch(currentBranch); err != nil {
			t.Errorf("Failed to restore original branch: %v", err)
		}
	}
}

// Note: The following tests are commented out as they modify the git repository
// They should be run in a test environment with a clean git repository

/*
func TestCreateStoryBranch(t *testing.T) {
	wm := NewWorkflowManager()
	storyID := "456"
	description := "Chat UI"

	err := wm.CreateStoryBranch(storyID, description)
	if err != nil {
		t.Errorf("Failed to create story branch: %v", err)
	}

	currentBranch, err := wm.GetCurrentBranch()
	if err != nil {
		t.Errorf("Failed to get current branch: %v", err)
	}

	expectedBranch := "W-456-chat-ui"
	if currentBranch != expectedBranch {
		t.Errorf("Expected branch %s, got %s", expectedBranch, currentBranch)
	}

	// Cleanup
	cmd := exec.Command("git", "checkout", "main")
	if err := cmd.Run(); err != nil {
		t.Errorf("Failed to cleanup: %v", err)
	}

	cmd = exec.Command("git", "branch", "-D", expectedBranch)
	if err := cmd.Run(); err != nil {
		t.Errorf("Failed to cleanup: %v", err)
	}
}

func TestCommitChanges(t *testing.T) {
	wm := NewWorkflowManager()
	scope := "chat"
	description := "add user message bubble"

	err := wm.CommitChanges(scope, description)
	if err != nil {
		t.Errorf("Failed to commit changes: %v", err)
	}

	// Verify the commit message
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	output, err := cmd.Output()
	if err != nil {
		t.Errorf("Failed to get commit message: %v", err)
	}

	expectedMessage := "feat(chat): add user message bubble"
	if strings.TrimSpace(string(output)) != expectedMessage {
		t.Errorf("Expected commit message %s, got %s", expectedMessage, string(output))
	}
}
*/
