# Hugo conversion log

## package main

| file | status | also see |
| ---- | ------ | -------- |
| `main.go` | done | pkg/commands/hugo.go |

Done.

- we don't call `runtime.GOMAXPROCS(runtime.NumCPU())` because this is the default as of Go 1.5
- we make `commands.Execute()` do the work of "if error log then error return"
- Bits of code moved to `core/init.go`.
- Hugo called os.Exit from inner functions, Drouet does not do that, to ensure that clean shutdown can be carried out.

All other package code is in the `pkg` directory.

## package commands

`hugo/commands/*` ==> `drouet/pkg/commands/*`

| hugo file | drouet file| Notes |
| ---- | ------ | -------- |
| benchmark.go         | benchmark.go         | cmd      |
| ~~check.go~~         | ~~check.go~~         | ~~done~~ |
| commandeer.go        |                      |          |
| convert.go           | convert.go           | cmd      |
| env.go               | env.go               | cmd      |
| gen.go               | gen.go               | cmd      |
| genautocomplete.go   | gen.go               | cmd      |
| genchromastyles.go   | gen.go               | cmd      |
| gendoc.go            | gen.go               | cmd      |
| gendocshelper.go     | gen.go               | cmd      |
| genman.go            | gen.go               | cmd      |
| hugo.go              | hugo.go              | cmd, config |
| ~~hugo_windows.go~~  | ~~hugo_windows.go~~  | ~~done~~ |
| import_jekyll.go     | import.go            | cmd      |
| import_jekyll_test.go|                      |          |
| ~~limit_darwin.go~~  | ~~check_darwin.go~~  | not tested |
| ~~limit_others.go~~ | ~~check_notdarwin.go~~| ~~done~~ |
| list.go              | list.go              | cmd      |
| list_config.go       | config.go            | except InitializeConfig |
| new.go               | new.go               | cmd      |
| new_test.go          |                      |          |
| release.go           | release.go           | except releaser |
| server.go            | server.go            | cmd      |
| ~~server_test.go~~   | ~~server_test.go~~   | ~~done~~ |
| static_syncer.go     |                      |          |
| ~~version.go~~       | ~~version.go~~       | ~~done~~ |

Each subcommand (verb) is moved to its own file. Some verbs have subcommands of their
own, and for now, each verb subcommand is in the verb's file.

The pattern for all command handlers is that there is a function that builds it and
returns it to be linked into a parent command. At the top-level, we build the root
command and execute it. Since there are no globals, we have a top-level variable that
holds necessary state for command execution; this is `struct core.Hugo` and is wrapped
inside a `struct HugoCmd` so we can use method variables to capture the receiver
in the `RunE` callback.

The root command handler is (still) in `hugo.go`. Maybe this should be named `root.go` per suggested
Cobra convention? It would be nice if this was sorted at the top of the list, but I rejected
the idea of naming files `hugo.go`, `hugo-benchmark.go` and so on, as that would result in
almost every file starting with the substring "hugo". There are no obvious words starting with
a legal filename character that would sort before the first verb, or stand out. It cannot
begin with "\_", for example. I could make it start with "0", I suppose.

Done

- hugo version

In process

- hugo
- hugo config

TBD

- hugo benchmark
- hugo convert
- hugo env
- hugo gen
- hugo import
- hugo list
- hugo new
- hugo server

### file hugo/commands/hugo.go

This is in progress; the command-line parsing and config loading is the first focus.

Done

- var HugoCmd \*cobra.Command is now a local
- `var hugoCmdV *cobra.Command` removed. ask why this was passed into flags parsing for all subcommands
- all global variables are removed
- Execute, AddCommands, initHugoBuilderFlags, initRootPersistentFlags
- initHugoBuildCommonFlags, initBenchmarkBuildingFlags
- init() function removed and code moved elsewhere
- initializeFlags replaced with "add all flags to Viper" (check for aliasing needed)

In progress

- InitializeConfig

TBD

- `var Hugo *hugolib.HugoSites` is evidently used by the Hugo caddy plugin. Talk to them?
- `Reset()`
- `commandError struct` and methods
- error handling not done yet
- createLogger
- deprecatedFlags
- fullBuild, build, serverBuild
- copyStatic, createStaticDirsConfig, doWithPublishDirs
- countingStatFs and methods
- copyStaticTo
- timeTrack
- getDirList
- recreateAndBuildSites, resetAndBuildSites, initSites, buildSites
- newWatcher
- pickOneWriteOrCreatePath, isThemeVsHugoVersionMismatch

#### commands.Execute()

Hugo and Drouet both start with `commands.Execute()`, which builds and runs the command-line
processor.

### file hugo/commands/list_config.go

`hugo/commands/list_config` ==> `drouet/pkg/commands/config.go`. I wonder if the Hugo name was
simply to avoid multiple files named `config.go`. I don't think this is worth doing.

The `hugo config` command got three new flags `--source`, `--theme`, and `--themesDir`, to
help locate config in cases where the site source and theme aren't colocated. This should have
been in Hugo to begin with, I'm guessing no one uses the `hugo config` command very much.
It should be enhanced a bit so that it can be used for better diagnosis of config setup.
Either that, or the `hugo check` command should exist and that would be one of the things
it does.

Done

- var configCmd \*cobra.Command is now a local
- printConfig

In progress

- InitializeConfig

### file hugo/commands/version.go

`hugo/commands/version.go` ==> `drouet/pkg/commands/version.go`.

This is a case where we kept a few global variables, so that the build process can write
values into them. The Hugo build is done like this

```
$ go install -X "github.com/gohugoio/hugo/hugolib.CommitHash=$COMMIT_HASH -X "github.com/gohugoio/hugo/hugolib.BuildDate=$BUILD_DATE"
```

which stamps values into `hugolib.CommitHash` and `hugolib.BuildDate`. These are now `core.CommitHash` and `core.BuildDate`.
Build by hand on Windows with

```
$ set XHASH=github.com/neurocline/drouet/pkg/core.CommitHash=B3587337
$ set XBUILD=github.com/neurocline/drouet/pkg/core.BuildDate=2018-03-31T13:13:28-0700
$ go install -v -ldflags "-X %XHASH% -X %XBUILD%" github.com/neurocline/drouet
```

Done

- var versionCmd \*cobra.Command is now a local
- printHugoVersion

TBD

- add mage file so that CommitHash and BuildDate are more easily stamped in

# By source file

- ~~main.go~~
- _bufferpool/_
- _cache/_
- commands/
  - benchmark.go
  - check.go
  - **commandeer.go**
  - convert.go
  - env.go
  - gen.go
  - genautocomplete.go
  - genchromastyles.go
  - gendoc.go
  - gendocshelper.go
  - genman.go
  - hugo.go
  - _hugo_windows.go_
  - import_jekyll.go
  - _import_jekyll_test.go_
  - _limit_darwin.go_
  - _limit_others.go_
  - list.go
  - ~~list_config.go~~
  - new.go
  - _new_test.go_
  - _release.go_
  - server.go
  - _server_test.go_
  - _static_syncer.go_
  - ~~version.go~~
