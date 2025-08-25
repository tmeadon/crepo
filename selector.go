package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

func findGitRepos(root, normalizedSearchTerm string, excludedDirs map[string]struct{}) ([]string, error) {
	var gitRepos []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			name := d.Name()
			if name == ".git" {
				repoPath := filepath.Dir(path)
				if isRepoMatch(repoPath, normalizedSearchTerm) {
					gitRepos = append(gitRepos, repoPath)
				}
				return filepath.SkipDir
			}

			if _, ok := excludedDirs[name]; ok {
				return filepath.SkipDir
			}
		}

		return nil
	})

	return gitRepos, err
}

func selectRepo(searchTerm string, config *Config) {
	normalizedSearch := normalize(searchTerm)

	repos, err := loadRepoCache()
	if err != nil || repos == nil {
		repos, err = findGitRepos(config.RepoRoot, normalizedSearch, config.exclusionSet)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error during repo discovery:", err)
			return
		}
		if err := saveRepoCache(repos); err != nil {
			fmt.Fprintln(os.Stderr, "Error saving repo cache:", err)
		}
	} else if searchTerm != "" {
		var filtered []string
		for _, repo := range repos {
			if isRepoMatch(repo, normalizedSearch) {
				filtered = append(filtered, repo)
			}
		}
		repos = filtered
	}

	switch len(repos) {
	case 0:
		fmt.Println("No repositories found")
		return
	case 1:
		switchDir(repos[0])
	default:
		selected := promptForSelection(repos)
		switchDir(selected)
	}
}

func switchDir(path string) {
	fmt.Print(path)
}

func isRepoMatch(path string, normalizedSearchTerm string) bool {
	return strings.Contains(normalize(path), normalizedSearchTerm)
}

func normalize(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "")
}

func promptForSelection(repos []string) string {
	searcher := func(input string, index int) bool {
		item := repos[index]
		return isRepoMatch(item, normalize(input))
	}

	prompt := promptui.Select{
		Label:             "Select a repository or press / to search",
		Items:             repos,
		Size:              10,
		Searcher:          searcher,
		StartInSearchMode: false,
		Stdout:            os.Stderr,
	}

	i, _, err := prompt.Run()

	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println("^C")
			os.Exit(0)
		}
		fmt.Println("Prompt failed:", err)
		return ""
	}

	return repos[i]
}
