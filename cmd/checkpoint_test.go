package cmd

import (
	"fmt"
	"strings"
	"testing"
)

// mockGitClient implements the GitClient interface without touching a real repo.
type mockGitClient struct {
	commitResult string
	commitErr    error
	latestHash   string
	changedFiles []string
}

func (m *mockGitClient) Commit(message string) (string, error) {
	return m.commitResult, m.commitErr
}

func (m *mockGitClient) LatestHash() (string, error) {
	return m.latestHash, nil
}

func (m *mockGitClient) ChangedLoadstarFiles() ([]string, error) {
	return m.changedFiles, nil
}

// setupCheckpointTest wires fs/svc AND replaces gitClient with the mock.
func setupCheckpointTest(t *testing.T, mock *mockGitClient) string {
	t.Helper()
	loadstarBase := setupCmdTest(t)
	gitClient = mock
	return loadstarBase
}

// ---- checkpoint ----

func TestCheckpoint_CommitSuccess(t *testing.T) {
	mock := &mockGitClient{commitResult: "abc123def456"}
	setupCheckpointTest(t, mock)

	hash, err := gitClient.Commit("test checkpoint")
	if err != nil {
		t.Fatalf("Commit: %v", err)
	}
	if hash != "abc123def456" {
		t.Errorf("expected abc123def456, got %s", hash)
	}
}

func TestCheckpoint_GitFailure(t *testing.T) {
	mock := &mockGitClient{commitErr: fmt.Errorf("nothing to commit")}
	setupCheckpointTest(t, mock)

	_, err := gitClient.Commit("fail")
	if err == nil {
		t.Error("expected error from failed commit")
	}
}

// ---- buildCheckpointMessage ----

func TestBuildCheckpointMessage_Basic(t *testing.T) {
	msg := buildCheckpointMessage("test msg", false, nil)
	if msg != "test msg" {
		t.Errorf("expected 'test msg', got %q", msg)
	}
}

func TestBuildCheckpointMessage_Auto(t *testing.T) {
	msg := buildCheckpointMessage("auto save", true, nil)
	if !strings.HasPrefix(msg, "[AUTO-CHECKPOINT] ") {
		t.Errorf("expected [AUTO-CHECKPOINT] prefix, got %q", msg)
	}
}

func TestBuildCheckpointMessage_WithChangedFiles(t *testing.T) {
	files := []string{
		".loadstar/WAYPOINT/root.cli.cmd_create.md",
		".loadstar/BLACKBOX/root.cli.cmd_create.md",
		".loadstar/.clionly/LOG/some.log.md",
	}
	msg := buildCheckpointMessage("test", false, files)
	if !strings.Contains(msg, "WAYPOINT/root.cli.cmd_create.md") {
		t.Error("should contain WAYPOINT file")
	}
	if !strings.Contains(msg, "BLACKBOX/root.cli.cmd_create.md") {
		t.Error("should contain BLACKBOX file")
	}
	if strings.Contains(msg, ".clionly") {
		t.Error("should not contain .clionly files")
	}
}
