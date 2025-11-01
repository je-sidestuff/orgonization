package templates

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
	Args  []DirectiveActionArg
}

type DirectiveActionArg int

const (
	MaxOccurrences DirectiveActionArg = iota + 1
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
