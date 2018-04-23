package banksy

import (
	"fmt"
	"testing"

	"github.com/google/go-github/github"
)

func TestGlobRuleSingleFileSingleGlob(t *testing.T) {
	f := github.CommitFile{}
	filename := "spring-boot-project/spring-boot/src/main/java/org/springframework/boot/util/LambdaSafe.java"
	f.Filename = &filename
	files := []*github.CommitFile{&f}

	gr := globRule{}
	gr.Globs = []interface{}{"*/boot/*"}

	match := gr.isMatch(nil, files)
	if !match {
		t.Error(fmt.Sprintf("Expected match, but returned false"))
	}
}

func TestGlobRuleSingleFileMultipleGlob(t *testing.T) {
	f := github.CommitFile{}
	filename := "spring-boot-project/spring-boot/src/main/java/org/springframework/boot/util/LambdaSafe.java"
	f.Filename = &filename
	files := []*github.CommitFile{&f}

	gr := globRule{}
	gr.Globs = []interface{}{"*/boot/*", "*foo*"}

	match := gr.isMatch(nil, files)
	if !match {
		t.Error(fmt.Sprintf("Expected match, but returned false"))
	}
}

func TestGlobRuleMulitpleFileSingleGlob(t *testing.T) {
	f := github.CommitFile{}
	filename := "spring-boot-project/spring-boot/src/main/java/org/springframework/boot/util/LambdaSafe.java"
	f.Filename = &filename
	f1 := github.CommitFile{}
	filename1 := "README.md"
	f1.Filename = &filename1
	files := []*github.CommitFile{&f, &f1}

	gr := globRule{}
	gr.Globs = []interface{}{"*/boot/*", "*foo*"}

	match := gr.isMatch(nil, files)
	if !match {
		t.Error(fmt.Sprintf("Expected match, but returned false"))
	}
}

func TestGlobRuleMulitpleFileSingleGlobNoMatch(t *testing.T) {
	f := github.CommitFile{}
	filename := "spring-boot-project/spring-boot/src/main/java/org/springframework/boot/util/LambdaSafe.java"
	f.Filename = &filename
	f1 := github.CommitFile{}
	filename1 := "README.md"
	f1.Filename = &filename1
	files := []*github.CommitFile{&f, &f1}

	gr := globRule{}
	gr.Globs = []interface{}{"*/bar/*", "*foo*"}

	match := gr.isMatch(nil, files)
	if match {
		t.Error(fmt.Sprintf("Expected no match, but returned true"))
	}
}

func TestSizeRuleGreaterThanNumFiles(t *testing.T) {
	sr := sizeRule{}
	sr.Compare = "greaterThan"
	sr.NumFiles = 3
	pr := github.PullRequest{}
	changedFiles := 5
	pr.ChangedFiles = &changedFiles

	match := sr.isMatch(&pr, []*github.CommitFile{})
	if !match {
		t.Error(fmt.Sprintf("Expected match, but returned false"))
	}
}

func TestSizeRuleGreaterThanNumFilesNoMatch(t *testing.T) {
	sr := sizeRule{}
	sr.Compare = "greaterThan"
	sr.NumFiles = 30
	pr := github.PullRequest{}
	changedFiles := 5
	pr.ChangedFiles = &changedFiles

	match := sr.isMatch(&pr, []*github.CommitFile{})
	if match {
		t.Error(fmt.Sprintf("Expected no match, but returned true"))
	}
}

func TestSizeRuleGreaterThanNumChanges(t *testing.T) {
	sr := sizeRule{}
	sr.Compare = "greaterThan"
	sr.NumChanges = 5
	pr := github.PullRequest{}
	additions := 2
	subtractions := 4
	pr.Additions = &additions
	pr.Deletions = &subtractions

	match := sr.isMatch(&pr, []*github.CommitFile{})
	if !match {
		t.Error(fmt.Sprintf("Expected match, but returned false"))
	}
}

func TestSizeRuleGreaterThanNumChangesNoMatch(t *testing.T) {
	sr := sizeRule{}
	sr.Compare = "greaterThan"
	sr.NumChanges = 50
	pr := github.PullRequest{}
	additions := 2
	subtractions := 4
	pr.Additions = &additions
	pr.Deletions = &subtractions

	match := sr.isMatch(&pr, []*github.CommitFile{})
	if match {
		t.Error(fmt.Sprintf("Expected no match, but returned true"))
	}
}

func TestSizeRuleGreaterThanWithChangesAndFilesGreater(t *testing.T) {
	sr := sizeRule{}
	sr.Compare = "greaterThan"
	sr.NumChanges = 5
	sr.NumFiles = 3
	pr := github.PullRequest{}
	additions := 2
	subtractions := 4
	pr.Additions = &additions
	pr.Deletions = &subtractions
	changedFiles := 6
	pr.ChangedFiles = &changedFiles

	match := sr.isMatch(&pr, []*github.CommitFile{})
	if !match {
		t.Error(fmt.Sprintf("Expected match, but returned false"))
	}
}

func TestSizeRuleGreaterThanWithChangesGreaterAndFilesLess(t *testing.T) {
	sr := sizeRule{}
	sr.Compare = "greaterThan"
	sr.NumChanges = 5
	sr.NumFiles = 5
	pr := github.PullRequest{}
	additions := 2
	subtractions := 4
	pr.Additions = &additions
	pr.Deletions = &subtractions
	changedFiles := 2
	pr.ChangedFiles = &changedFiles
	// when both NumChanges and NumFiles are specified it does an or
	match := sr.isMatch(&pr, []*github.CommitFile{})
	if !match {
		t.Error(fmt.Sprintf("Expected match, but returned false"))
	}
}

func TestSizeRuleGreaterThanWithChangesLessAndFilesGreater(t *testing.T) {
	sr := sizeRule{}
	sr.Compare = "greaterThan"
	sr.NumChanges = 50
	sr.NumFiles = 5
	pr := github.PullRequest{}
	additions := 2
	subtractions := 4
	pr.Additions = &additions
	pr.Deletions = &subtractions
	changedFiles := 10
	pr.ChangedFiles = &changedFiles
	// when both NumChanges and NumFiles are specified it does an or
	match := sr.isMatch(&pr, []*github.CommitFile{})
	if !match {
		t.Error(fmt.Sprintf("Expected match, but returned false"))
	}
}

/////
func TestSizeRuleLessThanNumFiles(t *testing.T) {
	sr := sizeRule{}
	sr.Compare = "LessThan"
	sr.NumFiles = 30
	pr := github.PullRequest{}
	changedFiles := 5
	pr.ChangedFiles = &changedFiles

	match := sr.isMatch(&pr, []*github.CommitFile{})
	if !match {
		t.Error(fmt.Sprintf("Expected match, but returned false"))
	}
}

func TestSizeRuleLessThanNumFilesNoMatch(t *testing.T) {
	sr := sizeRule{}
	sr.Compare = "lessThan"
	sr.NumFiles = 5
	pr := github.PullRequest{}
	changedFiles := 30
	pr.ChangedFiles = &changedFiles

	match := sr.isMatch(&pr, []*github.CommitFile{})
	if match {
		t.Error(fmt.Sprintf("Expected no match, but returned true"))
	}
}

func TestSizeRuleLessThanNumChanges(t *testing.T) {
	sr := sizeRule{}
	sr.Compare = "lessThan"
	sr.NumChanges = 50
	pr := github.PullRequest{}
	additions := 2
	subtractions := 4
	pr.Additions = &additions
	pr.Deletions = &subtractions

	match := sr.isMatch(&pr, []*github.CommitFile{})
	if !match {
		t.Error(fmt.Sprintf("Expected match, but returned false"))
	}
}

func TestSizeRuleLessThanNumChangesNoMatch(t *testing.T) {
	sr := sizeRule{}
	sr.Compare = "lessThan"
	sr.NumChanges = 50
	pr := github.PullRequest{}
	additions := 20
	subtractions := 40
	pr.Additions = &additions
	pr.Deletions = &subtractions

	match := sr.isMatch(&pr, []*github.CommitFile{})
	if match {
		t.Error(fmt.Sprintf("Expected no match, but returned true"))
	}
}

func TestSizeRuleLessThanWithChangesAndFilesLess(t *testing.T) {
	sr := sizeRule{}
	sr.Compare = "lessThan"
	sr.NumChanges = 5
	sr.NumFiles = 3
	pr := github.PullRequest{}
	additions := 2
	subtractions := 1
	pr.Additions = &additions
	pr.Deletions = &subtractions
	changedFiles := 1
	pr.ChangedFiles = &changedFiles

	match := sr.isMatch(&pr, []*github.CommitFile{})
	if !match {
		t.Error(fmt.Sprintf("Expected match, but returned false"))
	}
}

func TestSizeRuleLessThanWithChangesGreaterAndFilesLess(t *testing.T) {
	sr := sizeRule{}
	sr.Compare = "lessThan"
	sr.NumChanges = 5
	sr.NumFiles = 5
	pr := github.PullRequest{}
	additions := 2
	subtractions := 4
	pr.Additions = &additions
	pr.Deletions = &subtractions
	changedFiles := 2
	pr.ChangedFiles = &changedFiles
	// when both NumChanges and NumFiles are specified it does an or
	match := sr.isMatch(&pr, []*github.CommitFile{})
	if !match {
		t.Error(fmt.Sprintf("Expected match, but returned false"))
	}
}

func TestSizeRuleLessThanWithChangesLessAndFilesGreater(t *testing.T) {
	sr := sizeRule{}
	sr.Compare = "lessThan"
	sr.NumChanges = 50
	sr.NumFiles = 5
	pr := github.PullRequest{}
	additions := 2
	subtractions := 4
	pr.Additions = &additions
	pr.Deletions = &subtractions
	changedFiles := 10
	pr.ChangedFiles = &changedFiles
	// when both NumChanges and NumFiles are specified it does an or
	match := sr.isMatch(&pr, []*github.CommitFile{})
	if !match {
		t.Error(fmt.Sprintf("Expected match, but returned false"))
	}
}
