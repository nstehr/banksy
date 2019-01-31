package banksy

import (
	"reflect"
	"testing"
)

func TestNonMatchingRule(t *testing.T) {

	entry := RuleEntry{}
	entry["type"] = "madeUpType"
	entries := []RuleEntry{entry}

	r := buildRules(entries)

	if len(r) != 0 {
		t.Errorf("Expected no rules to be build, instead rule was built")
	}
}

func TestBuildGlobRule(t *testing.T) {
	entry := RuleEntry{}
	entry["type"] = "globRule"
	entry["Label"] = "globby"

	var globs []interface{}
	globs = append(globs, "first", "second")

	entry["Globs"] = globs

	entries := []RuleEntry{entry}

	r := buildRules(entries)
	if len(r) != 1 {
		t.Errorf("Expected 1 rules to be build, instead no rule was built")
	}
	gr := r[0]
	if reflect.TypeOf(gr) != reflect.TypeOf(&globRule{}) {
		t.Errorf("Unexpected type returned for rule[%v]: %v", "globRule", reflect.TypeOf(gr))
	}
	rule := gr.(*globRule)
	if rule.Label != "globby" {
		t.Errorf("Expected globby, received: %s", rule.Label)
	}
	if len(rule.Globs) != 2 {
		t.Errorf("Expected 2 entries, received: %d", len(rule.Globs))
	}
}

func TestBuildSizeRule(t *testing.T) {
	entry := RuleEntry{}
	entry["type"] = "sizeRule"
	entry["Compare"] = "GreaterThan"
	entry["NumFiles"] = 5
	entry["NumChanges"] = 10
	entry["Label"] = "Small"
	entries := []RuleEntry{entry}

	r := buildRules(entries)
	if len(r) != 1 {
		t.Errorf("Expected 1 rules to be build, instead no rule was built")
	}
	sr := r[0]
	if reflect.TypeOf(sr) != reflect.TypeOf(&sizeRule{}) {
		t.Errorf("Unexpected type returned for rule[%v]: %v", "sizeRule", reflect.TypeOf(sr))
	}

	rule := sr.(*sizeRule)
	if rule.Compare != "GreaterThan" {
		t.Errorf("Expected GreaterThan, received: %s", rule.Compare)
	}
	if rule.Label != "Small" {
		t.Errorf("Expected Small, received: %s", rule.Label)
	}
	if rule.NumFiles != 5 {
		t.Errorf("Expected 5, received: %d", rule.NumFiles)
	}
	if rule.NumChanges != 10 {
		t.Errorf("Expected 10, received: %d", rule.NumChanges)
	}
}
