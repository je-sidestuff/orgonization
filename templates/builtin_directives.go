package templates

import (
	"fmt"
	"time"
)

type builtinDirectiveTriggerClasses struct {
	matchCommandToken DirectiveTriggerClass
}

var builtinDirectiveTriggerClassesInstance = &builtinDirectiveTriggerClasses{
	matchCommandToken: DirectiveTriggerClass{
		Name: "MatchCommandToken",
		AllowedArgs: []DirectiveTriggerArg{
			Pattern,
		},
		Inputs: []DirectiveTriggerCondition{
			TokenMatched,
		},
		Outputs: []DirectiveTriggerTarget{
			MatchedToken,
		},
	},
}

type builtinDirectiveActionClasses struct {
	replaceTokenWithFuncResult DirectiveActionClass
}

var builtinDirectiveActionClassesInstance = &builtinDirectiveActionClasses{
	replaceTokenWithFuncResult: DirectiveActionClass{
		Name: "ReplaceTokenWithFuncResult",
		AllowedArgs: []DirectiveActionArg{
			PassArg,
			FunctionResult,
		},
		Targets: []DirectiveTriggerTarget{
			MatchedToken,
		},
		Effects: []DirectiveActionEffect{
			ReplaceToken,
		},
	},
}

func GetStatusResult(args []string) (string, error) {

	fmt.Println("Called GetStatusResult with args:", args)

	if len(args) < 1 {
		return "OK", nil
	}

	if args[0] == "TIME" {

		return time.Now().Format(time.RFC3339), nil
	}

	return args[0], nil
}

type builtinDirectives struct {
	weeklyNotesDirective DirectiveConfiguration
	statusDirective      DirectiveConfiguration
}

var builtinDirectivesInstance = &builtinDirectives{
	weeklyNotesDirective: DirectiveConfiguration{
		Name: "WeeklyNotesReplace",
		DirectiveActions: []DirectiveAction{
			{
				Name:  "UpdateWeeklyNotesToken",
				Class: builtinDirectiveActionClassesInstance.replaceTokenWithFuncResult,
				Args: map[DirectiveActionArg]string{
					PassArg:        "true",
					FunctionResult: `PrintWeekdays`,
				},
			},
		},
		DirectiveTriggers: []DirectiveTrigger{
			{
				Name:  "MatchWeeklyNotessToken",
				Class: builtinDirectiveTriggerClassesInstance.matchCommandToken,
				Args: map[DirectiveTriggerArg]string{
					Pattern: `<WEEKLY_NOTES:.{0,10}>`,
				},
			},
		},
	},
	statusDirective: DirectiveConfiguration{
		Name: "StatusCheck",
		DirectiveActions: []DirectiveAction{
			{
				Name:  "UpdateStatusToken",
				Class: builtinDirectiveActionClassesInstance.replaceTokenWithFuncResult,
				Args: map[DirectiveActionArg]string{
					PassArg:        "true",
					FunctionResult: `GetStatusResult`,
				},
			},
		},
		DirectiveTriggers: []DirectiveTrigger{
			{
				Name:  "MatchStatusToken",
				Class: builtinDirectiveTriggerClassesInstance.matchCommandToken,
				Args: map[DirectiveTriggerArg]string{
					Pattern: `<PING:.{0,10}>`,
				},
			},
		},
	},
}

func GetWeeklyTemplateDirective() DirectiveConfiguration {
	return builtinDirectivesInstance.weeklyNotesDirective
}

func GetDefaultTemplateDirectives() []DirectiveConfiguration {

	// Return just the weekly and status directives for now
	return []DirectiveConfiguration{
		builtinDirectivesInstance.weeklyNotesDirective,
		builtinDirectivesInstance.statusDirective}
}
