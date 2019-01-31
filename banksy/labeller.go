package banksy

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/google/go-github/github"
	yaml "gopkg.in/yaml.v2"
)

// RuleEntry represents a generic rule as read from yaml
type RuleEntry map[string]interface{}

// RuleEntrySet collection of RuleEntry maps
type RuleEntrySet struct {
	Entries []RuleEntry `yaml:"rules"`
}

// Labeller will label a PR based on a collection of rules
type Labeller struct {
	rules []Rule
}

// NewLabeller will instantiate a new Labeller
func NewLabeller(data []byte) (*Labeller, error) {
	var ruleSet RuleEntrySet
	err := yaml.Unmarshal([]byte(data), &ruleSet)
	if err != nil {
		return nil, err
	}
	rules := buildRules(ruleSet.Entries)
	return &Labeller{rules: rules}, nil
}

func buildRules(entries []RuleEntry) []Rule {
	var rules []Rule
	// I didn't know a way to unmarshal a yaml file containing structs of different types
	// every where I saw assumes the same type.  Instead of doing a custom unmarshalling
	// I deserialize to a map and use reflection to apply that map to the appropriate struct type
	for _, r := range entries {
		//TODO: come up with a cool rule registration mechanism
		var rule Rule
		switch r["type"] {
		case "globRule":
			rule = &globRule{}
		case "sizeRule":
			rule = &sizeRule{}
		default:
			log.Printf("no rule defined for type: %s\n", r["type"])
			continue
		}
		err := fillStruct(r, rule)
		if err != nil {
			log.Println("Error populating rule:", err)
		}
		rules = append(rules, rule)
	}
	return rules
}

// for a given PR apply the rules and label if criteria is met
func (l *Labeller) determineLabelling(pr *github.PullRequest, files []*github.CommitFile) []string {
	var labels []string
	for _, rule := range l.rules {
		if rule.isMatch(pr, files) {
			log.Printf("Found match for rule with label: %s\n", rule.getLabel())
			labels = append(labels, rule.getLabel())
		}
	}
	return labels
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)
	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		log.Println(val.Type())
		log.Println(structFieldType)
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func fillStruct(m map[string]interface{}, s interface{}) error {
	for k, v := range m {
		// just ignore the type, it will be in the yaml value, but in the rule structs
		if k == "type" {
			continue
		}
		err := setField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
