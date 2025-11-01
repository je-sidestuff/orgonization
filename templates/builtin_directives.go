package templates

type builtinDirectives struct {
	weeklyDirective DirectiveConfiguration
	statusDirective DirectiveConfiguration
}

var instance = &builtinDirectives{
	weeklyDirective: DirectiveConfiguration{
		Name: "WeeklyNotesGenerator",
		DirectiveActions: []DirectiveAction{
			{
				Name: "WeeklyNotesGenerator",
				Class: DirectiveActionClass{
					Name: "GenerateWeeklyNotes",
					AllowedArgs: []DirectiveActionArg{
						PassArg,
					},
					Targets: []DirectiveTriggerTarget{
						MatchedFile,
					},
					Effects: []DirectiveActionEffect{
						AppendLinesToFile,
					},
				},
				Args: []DirectiveActionArg{
					PassArg,
				},
			},
		},
		DirectiveTriggers: []DirectiveTrigger{
			{
				Name: "AlwaysTrigger",
				Class: DirectiveTriggerClass{
					Name: "AlwaysExecute",
					AllowedArgs: []DirectiveTriggerArg{
						AlwaysLog,
					},
					Inputs: []DirectiveTriggerCondition{
						Always,
					},
					Outputs: []DirectiveTriggerTarget{
						MatchedFile,
					},
				},
				Args: map[DirectiveTriggerArg]string{
					AlwaysLog: "true",
				},
			},
		},
	},
	statusDirective: DirectiveConfiguration{
		Name: "StatusCheck",
		DirectiveActions: []DirectiveAction{
			{
				Name: "UpdateStatusToken",
				Class: DirectiveActionClass{
					Name: "ReqRepTokenReplace",
					AllowedArgs: []DirectiveActionArg{
						PassArg,
					},
					Targets: []DirectiveTriggerTarget{
						MatchedToken,
					},
					Effects: []DirectiveActionEffect{
						ReplaceToken,
					},
				},
				Args: []DirectiveActionArg{
					PassArg,
				},
			},
		},
		DirectiveTriggers: []DirectiveTrigger{
			{
				Name: "MatchStatusToken",
				Class: DirectiveTriggerClass{
					Name: "ReqRepTokenMatch",
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
				Args: map[DirectiveTriggerArg]string{
					Pattern: `<PING:.{0,3}>`,
				},
			},
		},
	},
}

func GetWeeklyTemplateDirective() DirectiveConfiguration {
	return instance.weeklyDirective
}

func GetDefaultTemplateDirectives() []DirectiveConfiguration {

	// Return just the weekly and status directives for now
	return []DirectiveConfiguration{
		//instance.weeklyDirective,
		instance.statusDirective}
}
