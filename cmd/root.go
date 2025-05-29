package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"git-report-generator/internal/config"
	"git-report-generator/internal/generator"
	"git-report-generator/internal/git"

	"github.com/spf13/cobra"
)

var (
	repoPath    string
	dateFrom    string
	dateTo      string
	outputPath  string
	configPath  string
	authorEmail string
	branch      string
)

var rootCmd = &cobra.Command{
	Use:   "git-report-generator",
	Short: "Generate PDF reports of Git commits for a specified time period",
	Long: `Git Report Generator is a CLI tool that creates professional PDF reports
of Git commits for a specified author and time period. The reports include
commit details such as SHA, date, message, and description.

Example usage:
  git-report-generator --repo /path/to/repo --from 2024-01-01 --to 2024-01-31
  git-report-generator --repo . --from 2024-01-01 --to 2024-01-31 --author john@example.com`,
	RunE: runGenerate,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&repoPath, "repo", "r", ".", "Path to the Git repository")
	rootCmd.Flags().StringVarP(&dateFrom, "from", "f", "", "Start date (YYYY-MM-DD format)")
	rootCmd.Flags().StringVarP(&dateTo, "to", "t", "", "End date (YYYY-MM-DD format)")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output PDF file path (default: report_YYYY-MM-DD.pdf)")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
	rootCmd.Flags().StringVarP(&authorEmail, "author", "a", "", "Author email to filter commits (if empty, uses git config user.email)")
	rootCmd.Flags().StringVarP(&branch, "branch", "b", "", "Branch name to analyze (if empty, uses current branch)")

	rootCmd.MarkFlagRequired("from")
	rootCmd.MarkFlagRequired("to")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Validate and parse dates
	fromDate, err := time.Parse("2006-01-02", dateFrom)
	if err != nil {
		return fmt.Errorf("invalid from date format. Use YYYY-MM-DD: %w", err)
	}

	toDate, err := time.Parse("2006-01-02", dateTo)
	if err != nil {
		return fmt.Errorf("invalid to date format. Use YYYY-MM-DD: %w", err)
	}

	if fromDate.After(toDate) {
		return fmt.Errorf("from date cannot be after to date")
	}

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Get absolute path to repository
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for repository: %w", err)
	}

	// Initialize Git service
	gitService, err := git.NewService(absRepoPath)
	if err != nil {
		return fmt.Errorf("failed to initialize Git service: %w", err)
	}

	// Get author email if not provided
	if authorEmail == "" {
		authorEmail, err = gitService.GetUserEmail()
		if err != nil {
			return fmt.Errorf("failed to get user email from git config: %w", err)
		}
	}

	// Get branch name if not provided
	if branch == "" {
		branch, err = gitService.GetCurrentBranch()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
	}

	// Get repository name
	repoName := gitService.GetRepositoryName()

	// Get commits for the specified period and author
	commits, err := gitService.GetCommits(fromDate, toDate, authorEmail, branch)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	if len(commits) == 0 {
		fmt.Printf("No commits found for author %s between %s and %s on branch %s\n",
			authorEmail, dateFrom, dateTo, branch)
		return nil
	}

	// Generate output filename if not provided
	if outputPath == "" {
		outputPath = fmt.Sprintf("report_%s.pdf", time.Now().Format("2006-01-02"))
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Generate PDF report
	reportData := &generator.ReportData{
		Config:         cfg,
		RepositoryName: repoName,
		BranchName:     branch,
		AuthorEmail:    authorEmail,
		DateFrom:       fromDate,
		DateTo:         toDate,
		Commits:        commits,
	}

	pdfGenerator := generator.NewPDFGenerator()
	if err := pdfGenerator.Generate(reportData, outputPath); err != nil {
		return fmt.Errorf("failed to generate PDF report: %w", err)
	}

	fmt.Printf("✅ Report generated successfully: %s\n", outputPath)
	fmt.Printf("📊 Found %d commits for %s between %s and %s\n",
		len(commits), authorEmail, dateFrom, dateTo)

	return nil
}
