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

## Methodology

1. Fork `gohugoio/hugo` ==> `neurocline/drouet`.
2. the branch `hugo:drouet-translation` has code removed from it as is it added to `drouet`. This is needed because Drouet is not a line-for-line copy of Hugo.
3. `drouet-translation` is periodically rebased on top of `hugo:master`; conflicts show us changes in code that was already moved over.
4. A diary of changes is kept in [docs/Conversion.md](docs/Conversion.md).

## Code style

Get used to running `go fmt` on the code all the time. It has to be run before checkin,
and I should probably just wire it up so it's done on file save. Since `go fmt` only writes LF files,
we want all our Go source to be forced to LF line endings and hard tabs.

New Go projects should have these files in them from the start:

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

## History

This is the third attempt at a conversion. The first pass is in the branch `drouet:first-pass`
and the second in the branch `drouet:second-pass`.
