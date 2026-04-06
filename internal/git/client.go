package git

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// StatusInfo holds the git integration status for display.
type StatusInfo struct {
	RemoteURL        string
	Branch           string
	LatestHash       string
	UncommittedFiles int
}

// GetStatus returns the current git integration status.
// Non-fatal errors (no repo, no HEAD) result in partial info rather than failure.
func (c *Client) GetStatus() (StatusInfo, error) {
	cfg, _ := LoadConfig(c.loadstarBase)
	info := StatusInfo{
		RemoteURL: cfg.RemoteURL,
		Branch:    cfg.Branch,
	}

	repo, err := gogit.PlainOpen(c.repoPath)
	if err != nil {
		return info, nil
	}

	if ref, err := repo.Head(); err == nil {
		info.LatestHash = ref.Hash().String()
		info.Branch = ref.Name().Short()
	}

	if wt, err := repo.Worktree(); err == nil {
		if st, err := wt.Status(); err == nil {
			info.UncommittedFiles = len(st)
		}
	}

	return info, nil
}

// UnsetRemote removes the git_config.json and deletes the "origin" remote from the repo.
func (c *Client) UnsetRemote() error {
	cfgPath := configPath(c.loadstarBase)
	if err := os.Remove(cfgPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove git config: %w", err)
	}

	repo, err := gogit.PlainOpen(c.repoPath)
	if err != nil {
		return nil // no repo — nothing more to do
	}
	_ = repo.DeleteRemote(remoteName)
	return nil
}

const remoteName = "origin"

// Client implements internal.GitClient backed by go-git.
type Client struct {
	repoPath     string
	loadstarBase string // path to .loadstar/ directory
}

func NewClient(repoPath string) *Client {
	return &Client{
		repoPath:     repoPath,
		loadstarBase: filepath.Join(repoPath, ".loadstar"),
	}
}

// Init initialises a new git repository at repoPath if one does not already exist.
func (c *Client) Init() error {
	_, err := gogit.PlainOpen(c.repoPath)
	if err == nil {
		// Already a git repo — nothing to do.
		return nil
	}
	_, err = gogit.PlainInit(c.repoPath, false)
	if err != nil {
		return fmt.Errorf("git init: %w", err)
	}
	return nil
}

// SetRemote saves remote URL, branch, author info, and PAT to git_config.json
// and registers the remote in the git repository.
func (c *Client) SetRemote(remoteURL, branch, userName, userEmail, pat string) error {
	cfg := Config{
		RemoteURL: remoteURL,
		Branch:    branch,
		UserName:  userName,
		UserEmail: userEmail,
		PAT:       pat,
	}
	if err := SaveConfig(c.loadstarBase, cfg); err != nil {
		return err
	}

	repo, err := gogit.PlainOpen(c.repoPath)
	if err != nil {
		return fmt.Errorf("open repo: %w", err)
	}

	// Remove existing remote if present, then re-add.
	_ = repo.DeleteRemote(remoteName)
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{remoteURL},
	})
	if err != nil {
		return fmt.Errorf("set remote: %w", err)
	}
	return nil
}

// ChangedLoadstarFiles returns a list of .loadstar/ files that have uncommitted changes.
// Each path is relative to the repo root (e.g. ".loadstar/WAYPOINT/root.cli.cmd_create.md").
func (c *Client) ChangedLoadstarFiles() ([]string, error) {
	repo, err := gogit.PlainOpen(c.repoPath)
	if err != nil {
		return nil, fmt.Errorf("open repo: %w", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("worktree: %w", err)
	}

	st, err := wt.Status()
	if err != nil {
		return nil, fmt.Errorf("status: %w", err)
	}

	var files []string
	for path, s := range st {
		if s.Worktree == gogit.Unmodified && s.Staging == gogit.Unmodified {
			continue
		}
		if len(path) > len(".loadstar/") && path[:len(".loadstar/")] == ".loadstar/" {
			files = append(files, path)
		}
	}
	return files, nil
}

// Commit stages all changes under .loadstar/ and creates a commit.
// Returns the resulting commit hash string.
func (c *Client) Commit(message string) (string, error) {
	repo, err := gogit.PlainOpen(c.repoPath)
	if err != nil {
		return "", fmt.Errorf("open repo: %w", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("worktree: %w", err)
	}

	// Stage all changes under .loadstar/
	loadstarGlob := filepath.Join(".loadstar", "*")
	if err := wt.AddGlob(loadstarGlob); err != nil {
		return "", fmt.Errorf("stage .loadstar/: %w", err)
	}

	// Use author from config if available, fallback to defaults.
	cfg, _ := LoadConfig(c.loadstarBase)
	authorName := cfg.UserName
	authorEmail := cfg.UserEmail
	if authorName == "" {
		authorName = "loadstar"
	}
	if authorEmail == "" {
		authorEmail = "loadstar@local"
	}

	hash, err := wt.Commit(message, &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return "", fmt.Errorf("commit: %w", err)
	}

	return hash.String(), nil
}

// Push pushes the current branch to origin using the PAT from LOADSTAR_GIT_TOKEN.
// Returns an error if no remote is configured or the token is missing.
func (c *Client) Push() error {
	cfg, err := LoadConfig(c.loadstarBase)
	if err != nil {
		return err
	}
	if cfg.RemoteURL == "" {
		return fmt.Errorf("no remote configured — run `loadstar init --remote <URL>` first")
	}

	token := Token(cfg)
	if token == "" {
		return fmt.Errorf("LOADSTAR_GIT_TOKEN environment variable is not set")
	}

	repo, err := gogit.PlainOpen(c.repoPath)
	if err != nil {
		return fmt.Errorf("open repo: %w", err)
	}

	branch := cfg.Branch
	if branch == "" {
		branch = "main"
	}
	refSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch)

	err = repo.Push(&gogit.PushOptions{
		RemoteName: remoteName,
		RefSpecs:   []config.RefSpec{config.RefSpec(refSpec)},
		Auth: &http.BasicAuth{
			Username: cfg.UserName,
			Password: token,
		},
	})
	if err == gogit.NoErrAlreadyUpToDate {
		return nil
	}
	if err != nil {
		return fmt.Errorf("push: %w", err)
	}
	return nil
}

// LatestHash returns the HEAD commit hash of the repository.
func (c *Client) LatestHash() (string, error) {
	repo, err := gogit.PlainOpen(c.repoPath)
	if err != nil {
		return "", fmt.Errorf("open repo: %w", err)
	}

	ref, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("HEAD: %w", err)
	}

	return ref.Hash().String(), nil
}
