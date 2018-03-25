![Drouet](https://raw.githubusercontent.com/neurocline/drouet/master/docs/drouet.png)

Drouet is an experimental rewrite of [Hugo](https://github.com/gohugoio/hugo), "A Fast and Flexible Static Site Generator built with love". You don't want to use this, you want to use Hugo instead.

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
