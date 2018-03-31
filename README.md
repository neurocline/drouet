[![Drouet](https://raw.githubusercontent.com/neurocline/drouet/master/docs/drouet.png)](https://en.wikipedia.org/wiki/Juliette_Drouet)

Drouet is an experimental rewrite of [Hugo](https://github.com/gohugoio/hugo), "A Fast and Flexible Static Site Generator built with love". You don't want to use this, you want to use Hugo instead.

[Go report card](https://goreportcard.com/report/github.com/neurocline/drouet)

## Overview

I love Hugo, but I don't love the current structure of the code. This is an experiment,
or really, a set of experiments.

In particular, this is still Hugo. All the code says "Hugo", copyright notices are left intact
where the code is partially or entirely from the Hugo project, and the goal here is not to
replace Hugo but to try out some radical refactorings.

## Goals

No globals. No global state. No hidden dependencies. Test code doesn't need huge amounts
of setup to test functionality. Easier performance monitoring.

Hugo has a large amount of tests, and this is good, but the tests themselves are quite
large or very dependent on incidental details, and this is not so good.

### No global state

This means no `init()` functions, no global variables.

### Errors

Test code should not need to string-match against `error` string values; that is, the interpretation
of `error` values by code should be through something other than what `Error()` is returning.

## Code style

Get used to running `go fmt` on the code all the time. It has to be run before checkin,
and I should probably just wire it up so it's done on file save. It's totally unfortunate
that `go fmt` only writes LF files, this just doesn't work well with typical Windows Git
behavior. Maybe this means all Go-only Git projects should be checked out with LF line endings
even on Windows?

OK, so I bit the bullet, and new Go projects should have these files in them from the start:

`.gitattributes` file that makes .go files checked out with LF line endings, because the `go fmt`
tool writes them that way.

```
# Text files have auto line endings
* text = auto

# Go source files always have LF line endings
*.go text eol=lf
```

`.editorconfig` that makes .go files have hard tabs in them and LF line endings, because the `go
fmt` tool writes them that way. It also says indent_size = 8, but we can ignore that part and we
want our files viewed with 4-space tabs.

```
root = true

[*]
indent_style = space
indent_size = 4
tab_width = 4
trim_trailing_whitespace = true
insert_final_newline = true

[*.go]
indent_style = tab
end_of_line = lf
```

`.gitignore` that ignores `vendor/` because we eventually will use `go vendor` or `mage`.

```
/vendor/
```

## Documentation

The main README.md should always be a minimal and sufficient guide to using the project.
Don't take the current state of this README.md as the example - right now, this contains
notes about work-in-progress, and all of that will be moved to a better place "soon".

It's too bad that there are so many files that want to go in the root of the project. The
Hugo project has 20 files and directories that need to be in the root, which pushes down
the readable documentation off the screen. I introduced a top-level `pkg` directory to
hold all Hugo source, but this won't be enough.
