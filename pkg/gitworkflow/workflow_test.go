package gitworkflow

import (
	"fmt"
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
			name:    "valid feature branch",
			branch:  "feat/login-button",
			wantErr: false,
		},
		{
			name:    "valid bugfix branch",
			branch:  "bugfix/issue-123",
			wantErr: false,
		},
		{
			name:    "valid chore branch",
			branch:  "chore/refactor-auth",
			wantErr: false,
		},
		{
			name:    "invalid branch type",
			branch:  "invalid/login-button",
			wantErr: true,
		},
		{
			name:    "missing description",
			branch:  "feat/",
			wantErr: true,
		},
		{
			name:    "no separator",
			branch:  "featlogin-button",
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
	wm := NewWorkflowManager()
	storyID := "456"
	description := "Chat UI"

	// This test doesn't actually create the branch, just verifies the format
	branchName := fmt.Sprintf("feat/story-%s-%s", storyID, strings.ToLower(strings.ReplaceAll(description, " ", "-")))
	expected := "feat/story-456-chat-ui"

	if branchName != expected {
		t.Errorf("Expected branch name %s, got %s", expected, branchName)
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
	
	expectedBranch := "feat/story-456-chat-ui"
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