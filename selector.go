package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

func findGitRepos(root, searchTerm string, excludedDirs []string) ([]string, error) {
	var gitRepos []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == ".git" {
			repoPath := filepath.Dir(path)
			if isRepoMatch(repoPath, searchTerm) {
				gitRepos = append(gitRepos, repoPath)
			}
			return filepath.SkipDir
		}

		// Skip directories in skipDirs
		for _, dir := range excludedDirs {
			if info.IsDir() && info.Name() == dir {
				return filepath.SkipDir
			}
		}

		return nil
	})

	return gitRepos, err
}

func selectRepo(searchTerm string, config *Config) {
	repos, err := findGitRepos(config.RepoRoot, searchTerm, config.Exclusions)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error during repo discovery:", err)
		return
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

func isRepoMatch(path string, searchTerm string) bool {
	path = strings.Replace(strings.ToLower(path), " ", "", -1)
	searchTerm = strings.Replace(strings.ToLower(searchTerm), " ", "", -1)
	return strings.Contains(path, searchTerm)
}

func promptForSelection(repos []string) string {
	searcher := func(input string, index int) bool {
		item := repos[index]
		return isRepoMatch(item, input)
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
