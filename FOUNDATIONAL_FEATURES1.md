# Foundational Features 1

This is a targets list for the items to include in the first foundational features bundle.

We'll include:
- Some cleanup
- Beginnings of the 'server' mode
- Refactoring of the 'builtin' templates to use the standard templating path
- Inclusion of an expandable token to drive weekly notes
- Finish ping-pong exchange (Status OK and the time or similar?) (good)
- README update to explain basic usage
  - "Test Drive" next increment

## Notes for cleanup

For cleanup first we'll go file-by-file and comment on things that need to be adjusted, with each item having a (now) or (later) speculation included.

After that we'll mention any other more abstract notes about cleaning up which are not directly linked to the code.

### File-by-file

./templates/weekly_notes_test.go
  - Fine for now, likely to expand as we add hooks to printweekdays

./templates/directive_configuration.go
  - Make just enough changes here to get through the immediate use-case (as use-cases are done)
  - This will be refactored more thoroughly when directive package 1 is added (later)

./templates/file_processor.go
  - Get rid of Name stuff or implement it (now)

./templates/filesystem_configuration.go
  - Get rid of Name stuff or make it useful (now)

./templates/weekly_notes.go
  - Refactor to output strings rather than printing (now) (good)

./templates/template_processor.go
  - Fix warnings! (now) (good)
  - Adjust mapping of tokens to replace so we can go back-to-front on the full list, not per-type (later)
  - Get rid of Name stuff or make it useful (now)
  - Some comments to fix (now)
  - Swap out the hardcoded replacer-string (now)

./templates/templates.go
  - Is good

./templates/builtin_directives.go
  - re-enable weekly directive (now)
  - adjust weekly so that it can be replaced into file (now)
  - minimize changes for now, bigger refactor comes with directive package 1 (later)

./integration/status_template_test.go
  - Adjust so it's using filesystem configuration the "normal way" (now)
  - Confirm specific content in output (now)

./integration/integration.go
  - Is good

./integration/filesystem_integration_test_manager.go
  - Fine for now

./integration/helpers.go
  - Is good

./integration/filesystem_integration_test_manager_test.go
  - Some slight corrections to update for newer project folder name (now)

./cmd/main.go
  - Make entrypoints more similar between builtin, once, and agent
  - Trim extra arguments/paths that are not used/needed yet (now)

### More General

- AlwaysTrigger
  - Will this actually make sense to have?

- Verify the debugger is set up, ensure included in devcontainer
  - document usage

## Beginnings of the 'server' mode

## Refactoring of the 'builtin' templates to use the standard templating path

- This means we need to define a 'special file' for STDIN/OUT or similar

## Inclusion of an expandable token to drive weekly notes

## Finish ping-pong exchange (Status OK and the time or similar?)

## README update to explain basic usage
