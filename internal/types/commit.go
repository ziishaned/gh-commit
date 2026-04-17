package types

// GitChange represents a single file change in git
type GitChange struct {
	Path      string // File path
	Additions int    // Lines added
	Deletions int    // Lines deleted
	Patch     string // Full diff patch
	IsNewFile bool   // Is this a new file
	IsDeleted bool   // Is this a deleted file
	IsRenamed bool   // Is this a renamed file
}

// CommitResult represents the result of a commit operation
type CommitResult struct {
	Success bool   // Was the commit successful
	Message string // Generated commit message
	Hash    string // Commit hash (if successful)
	Error   error  // Error if unsuccessful
	DryRun  bool   // Was this a dry run?
}

// LLMRequest represents a request to the LLM
type LLMRequest struct {
	Diff     string // Git diff to analyze
	Model    string // Model to use
	Language string // Language for output (unused but kept for compatibility)
}

// LLMResponse represents a response from the LLM
type LLMResponse struct {
	Message string // Generated commit message
	Error   error  // Error if unsuccessful
}
