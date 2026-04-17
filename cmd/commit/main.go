package main

import (
	"fmt"
	"os"

	"github.com/gh-commit/internal/git"
	"github.com/spf13/cobra"
)

const extensionName = "commit"

var rootCmd = &cobra.Command{
	Use:   extensionName,
	Short: "Generate AI-powered commit messages",
	Long:  "A GitHub CLI extension that generates commit messages using GitHub Models and commits them automatically",
	RunE:  runCommit,
}

var (
	flagModel  string
	flagDryRun bool
)

func init() {
	rootCmd.Flags().StringVarP(&flagModel, "model", "m", "openai/gpt-4o", "GitHub Models model to use")
	rootCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Show what would be done without actually committing")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runCommit(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking for staged changes...")

	gitClient := &git.Client{}

	// Check if there are any changes at all
	if !gitClient.HasAnyChanges() {
		return fmt.Errorf("no changes detected. Please make some changes first")
	}

	// Check if there are staged changes
	if !gitClient.HasStagedChanges() {
		fmt.Println("No staged changes found. Staging all changes...")
		if err := gitClient.StageAllChanges(); err != nil {
			return fmt.Errorf("failed to stage changes: %w", err)
		}
		fmt.Println("✅ Changes staged")
	} else {
		fmt.Println("✅ Found staged changes")
	}

	// Get the staged diff
	fmt.Println("Analyzing staged changes...")
	diff, err := gitClient.GetStagedDiff()
	if err != nil {
		return fmt.Errorf("failed to get diff: %w", err)
	}

	fmt.Printf("📝 Analyzing %d lines of changes...\n", len(diff))

	return nil
}
