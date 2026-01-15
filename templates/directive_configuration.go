package templates

import (
	"fmt"
	"strings"
	"time"
)

// A DirectiveConfiguration defines the set of directives to be loaded for a given
// template processor.
// It contains the triggers which will cause a directive to execute and the actions
// which will be performed when the directive executes.
type DirectiveConfiguration struct {
	Name              string
	DirectiveActions  []DirectiveAction
	DirectiveTriggers []DirectiveTrigger
	// Will we want to supply args directly to the directive?
}

type DirectiveActionClass struct {
	Name        string
	AllowedArgs []DirectiveActionArg
	Targets     []DirectiveTriggerTarget
	Effects     []DirectiveActionEffect
}

type DirectiveAction struct {
	Name  string
	Class DirectiveActionClass
	Args  map[DirectiveActionArg]string
}

func (d DirectiveAction) GetReplacementString(mappedToken string) (string, error) {

	// If the directive action class is ReplaceTokenWithFuncResult then
	// we will call the replacement function with the args and return the result
	if d.Class.Name == "ReplaceTokenWithFuncResult" {

		// Switch on the arg FunctionResult to decide what funmction to call
		if d.Args[FunctionResult] == "GetStatusResult" {

			// We have a string in the form of <CMD:ARG> -- if ARG is empty then we pass an empty slice
			// to the replacement function
			afterColon := strings.Split(mappedToken, ":")[1]
			beforeClose := strings.Split(afterColon, ">")[0]

			var args []string
			if beforeClose != "" {
				args = []string{beforeClose}
			} else {
				args = []string{}
			}

			// Get the status result with a slice containing only mappedToken
			result, err := GetStatusResult(args)

			if err != nil {
				return "", err
			}

			return "<PING:" + result + ">", nil
		}
		if d.Args[FunctionResult] == "PrintWeekdays" {
			// We have a string in the form of <CMD:ARG> -- if ARG is empty then we pass an empty slice
			// to the replacement function
			afterColon := strings.Split(mappedToken, ":")[1]
			beforeClose := strings.Split(afterColon, ">")[0]

			dateArg := time.Now().AddDate(0, 0, 7) // Default to one week from today
			if parsedDate, err := time.Parse("2006-01-02", beforeClose); err == nil {
				dateArg = parsedDate
			}

			// Get the status result with a slice containing only mappedToken
			result, err := PrintWeekdays(dateArg)

			if err != nil {
				return "", err
			}

			return result, nil
		}
	}

	return "", nil
}

type DirectiveActionArg int

const (
	MaxOccurrences DirectiveActionArg = iota + 1
	FunctionResult
	PassArg
)

type DirectiveActionEffect int

const (
	ReplaceToken DirectiveActionEffect = iota + 1
	InsertLinesAfter
	AppendLinesToFile
)

type DirectiveTriggerClass struct {
	Name        string
	AllowedArgs []DirectiveTriggerArg
	Inputs      []DirectiveTriggerCondition
	Outputs     []DirectiveTriggerTarget
}

type DirectiveTrigger struct {
	Name  string
	Class DirectiveTriggerClass
	Args  map[DirectiveTriggerArg]string
}

type DirectiveTriggerArg int

const (
	OnlyOnTick DirectiveTriggerArg = iota + 1
	AlwaysLog
	Pattern
)

type DirectiveTriggerCondition int

const (
	Always DirectiveTriggerCondition = iota + 1
	FileExists
	FileChanged
	FileCreated
	TokenMatched
	SectionContains
)

type DirectiveTriggerTarget int

const (
	MatchedToken DirectiveTriggerTarget = iota + 1
	MatchedFile
)

// String returns a nicely formatted multi-line string representation of the DirectiveConfiguration
func (dc DirectiveConfiguration) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Directive: %s\n", dc.Name))

	// Format actions
	sb.WriteString("  Actions:\n")
	for _, action := range dc.DirectiveActions {
		sb.WriteString(fmt.Sprintf("    - %s (Class: %s)\n", action.Name, action.Class.Name))
		sb.WriteString("      Targets: [")
		for i, target := range action.Class.Targets {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(targetToString(target))
		}
		sb.WriteString("]\n")
		sb.WriteString("      Effects: [")
		for i, effect := range action.Class.Effects {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(effectToString(effect))
		}
		sb.WriteString("]\n")
		if len(action.Args) > 0 {
			sb.WriteString("      Args:\n")
			for arg, value := range action.Args {
				sb.WriteString(fmt.Sprintf("        %s: %s\n", argToString(arg), value))
			}
		}
	}

	// Format triggers
	sb.WriteString("  Triggers:\n")
	for _, trigger := range dc.DirectiveTriggers {
		sb.WriteString(fmt.Sprintf("    - %s (Class: %s)\n", trigger.Name, trigger.Class.Name))
		sb.WriteString("      Conditions: [")
		for i, condition := range trigger.Class.Inputs {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(conditionToString(condition))
		}
		sb.WriteString("]\n")
		if len(trigger.Args) > 0 {
			sb.WriteString("      Args:\n")
			for arg, value := range trigger.Args {
				sb.WriteString(fmt.Sprintf("        %s: %s\n", triggerArgToString(arg), value))
			}
		}
	}

	return sb.String()
}

// Helper functions to convert enums to strings
func targetToString(target DirectiveTriggerTarget) string {
	switch target {
	case MatchedToken:
		return "MatchedToken"
	case MatchedFile:
		return "MatchedFile"
	default:
		return fmt.Sprintf("Unknown(%d)", target)
	}
}

func effectToString(effect DirectiveActionEffect) string {
	switch effect {
	case ReplaceToken:
		return "ReplaceToken"
	case InsertLinesAfter:
		return "InsertLinesAfter"
	case AppendLinesToFile:
		return "AppendLinesToFile"
	default:
		return fmt.Sprintf("Unknown(%d)", effect)
	}
}

func argToString(arg DirectiveActionArg) string {
	switch arg {
	case MaxOccurrences:
		return "MaxOccurrences"
	case FunctionResult:
		return "FunctionResult"
	case PassArg:
		return "PassArg"
	default:
		return fmt.Sprintf("Unknown(%d)", arg)
	}
}

func triggerArgToString(arg DirectiveTriggerArg) string {
	switch arg {
	case OnlyOnTick:
		return "OnlyOnTick"
	case AlwaysLog:
		return "AlwaysLog"
	case Pattern:
		return "Pattern"
	default:
		return fmt.Sprintf("Unknown(%d)", arg)
	}
}

func conditionToString(condition DirectiveTriggerCondition) string {
	switch condition {
	case Always:
		return "Always"
	case FileExists:
		return "FileExists"
	case FileChanged:
		return "FileChanged"
	case FileCreated:
		return "FileCreated"
	case TokenMatched:
		return "TokenMatched"
	case SectionContains:
		return "SectionContains"
	default:
		return fmt.Sprintf("Unknown(%d)", condition)
	}
}
