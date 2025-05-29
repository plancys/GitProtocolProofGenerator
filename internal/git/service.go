package git

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Commit represents a Git commit with relevant information
type Commit struct {
	SHA         string
	Date        time.Time
	Message     string
	Description string
	Author      string
	AuthorEmail string
}

// Service provides Git repository operations
type Service struct {
	repo     *git.Repository
	repoPath string
}

// NewService creates a new Git service for the specified repository path
func NewService(repoPath string) (*Service, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Git repository at %s: %w", repoPath, err)
	}

	return &Service{
		repo:     repo,
		repoPath: repoPath,
	}, nil
}

// GetRepositoryName returns the name of the repository
func (s *Service) GetRepositoryName() string {
	return filepath.Base(s.repoPath)
}

// GetCurrentBranch returns the current branch name
func (s *Service) GetCurrentBranch() (string, error) {
	head, err := s.repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	branchName := head.Name().Short()
	return branchName, nil
}

// GetUserEmail returns the user email from Git configuration
func (s *Service) GetUserEmail() (string, error) {
	config, err := s.repo.Config()
	if err != nil {
		return "", fmt.Errorf("failed to get repository config: %w", err)
	}

	if config.User.Email == "" {
		return "", fmt.Errorf("user email not configured in Git")
	}

	return config.User.Email, nil
}

// GetCommits retrieves commits for the specified author, date range, and branch
func (s *Service) GetCommits(fromDate, toDate time.Time, authorEmail, branchName string) ([]*Commit, error) {
	// Get the branch reference
	branchRefName := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branchName))
	branchRef, err := s.repo.Reference(branchRefName, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch reference for %s: %w", branchName, err)
	}

	// Get commit iterator
	commitIter, err := s.repo.Log(&git.LogOptions{
		From: branchRef.Hash(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}
	defer commitIter.Close()

	var commits []*Commit

	// Iterate through commits
	err = commitIter.ForEach(func(c *object.Commit) error {
		// Check if commit is within date range
		if c.Author.When.Before(fromDate) || c.Author.When.After(toDate.Add(24*time.Hour)) {
			return nil
		}

		// Check if commit is by the specified author
		if !strings.EqualFold(c.Author.Email, authorEmail) {
			return nil
		}

		// Parse commit message and description
		message, description := parseCommitMessage(c.Message)

		commit := &Commit{
			SHA:         c.Hash.String()[:8], // Short SHA
			Date:        c.Author.When,
			Message:     message,
			Description: description,
			Author:      c.Author.Name,
			AuthorEmail: c.Author.Email,
		}

		commits = append(commits, commit)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate through commits: %w", err)
	}

	// Sort commits by date (newest first)
	for i := 0; i < len(commits)-1; i++ {
		for j := i + 1; j < len(commits); j++ {
			if commits[i].Date.Before(commits[j].Date) {
				commits[i], commits[j] = commits[j], commits[i]
			}
		}
	}

	return commits, nil
}

// parseCommitMessage separates the commit message into title and description
func parseCommitMessage(fullMessage string) (message, description string) {
	lines := strings.Split(strings.TrimSpace(fullMessage), "\n")

	if len(lines) == 0 {
		return "", ""
	}

	message = strings.TrimSpace(lines[0])

	if len(lines) > 1 {
		// Join remaining lines as description, skipping empty lines
		var descLines []string
		for i := 1; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if line != "" {
				descLines = append(descLines, line)
			}
		}
		description = strings.Join(descLines, " ")
	}

	return message, description
}
