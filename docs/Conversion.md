# Hugo conversion log

## Methodology

1. Fork `gohugoio/hugo` ==> `neurocline/drouet`.
2. the branch `hugo:drouet-translation` has code removed from it as is it added to `drouet`. This is needed because Drouet is not a line-for-line copy of Hugo.
3. `drouet-translation` is periodically rebased on top of `hugo:master`; conflicts show us changes in code that was already moved over.

## Big changes

Umm, I need to make config loading reentrant. Maybe I won't confuse the code for now and leave
`commands.commandeer` in place, just hoisted to the top level so there are no globals.

## Changelog

- moved `main.go`.
- moved `commands.Execute()`, `Hugo` global, top-level `HugoCmd`, and some small related code like versioning and the build globals.
- moved "hugo benchmark", "hugo check", "hugo config", "hugo convert", "hugo env", "hugo list", "hugo version", as well as stub for config and parser
- moved "hugo gen", "hugo import", "hugo new", "hugo release", "hugo server", and stub code
- Added IntializeConfig and createLogger
