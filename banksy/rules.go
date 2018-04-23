package banksy

import (
	"fmt"
	"log"
	"strings"

	"github.com/gobwas/glob"
	"github.com/google/go-github/github"
)

const (
	greaterThan = "greaterthan"
	lessThan    = "lessthan"
)

// Rule is represents a rule to run a pull request against
type Rule interface {
	isMatch(pr *github.PullRequest, files []*github.CommitFile) bool
	getLabel() string
}

// a rule that labels a PR based on whether it's changed files match a glob from
// the collection of globs
type globRule struct {
	Label string
	//TODO: see if I can it to be slice of string instead of interface{} limited by my
	// building of these structs using reflection
	Globs []interface{}
}

// isMatch will return a match if one or more of the files changed matches
// any of the globs
func (gr *globRule) isMatch(pr *github.PullRequest, files []*github.CommitFile) bool {
	for _, globStr := range gr.Globs {
		gs, ok := globStr.(string)
		if !ok {
			fmt.Printf("Not a string, Value:'%v'\n", globStr)
			continue
		}
		g := glob.MustCompile(gs)
		for _, file := range files {
			if g.Match(file.GetFilename()) {
				return true
			}
		}
	}

	return false
}

func (gr *globRule) getLabel() string {
	return gr.Label
}

type sizeRule struct {
	Compare    string
	NumFiles   int
	NumChanges int
	Label      string
}

func (sr *sizeRule) getLabel() string {
	return sr.Label
}

func (sr *sizeRule) isMatch(pr *github.PullRequest, files []*github.CommitFile) bool {
	compare := strings.TrimSpace(sr.Compare)
	compare = strings.ToLower(compare)
	if compare == "" {
		log.Println("Invalid rule, must specify a Compare field")
		return false
	}
	if compare != greaterThan && compare != lessThan {
		log.Println("Invalid rule, must specify greaterThan or lessThan")
		return false
	}

	if compare == greaterThan {
		if (pr.GetChangedFiles() > sr.NumFiles) && sr.NumFiles > 0 ||
			((pr.GetAdditions()+pr.GetDeletions()) > sr.NumChanges) && sr.NumChanges > 0 {
			return true
		}
	}

	if compare == lessThan {
		if pr.GetChangedFiles() < sr.NumFiles ||
			(pr.GetAdditions()+pr.GetDeletions()) < sr.NumChanges {
			return true
		}
	}

	return false
}
