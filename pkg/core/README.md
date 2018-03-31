The core package is what hugolib used to be called.

## HugoVersion

The Hugo code has what seems like an overly-complicated set of methods to declare and
manipulate versions. What does the code outside `helpers/hugo_info.go` see?

- `commands/commandeer.go`
  - `commands.loadConfig()` calls `helpers.CurrentHugoVersion.ReleaseVersion()` to add the release version to an error message.
- `commands/genman.go`
  - `commands.genmanCmd.RunE` calls `helpers.CurrentHugoVersion.String()` to put the full Hugo version into the header of each generated doc file.
- `commands/version.go`
  - `commands.printHugoVersion()` calls `helpers.CurrentHugoVersion.String()` to show the Hugo version to the user.
- `commands/hugo.go`
  - `commands.isThemeVsHugoVersionMismatch` calls `helpers.CompareVersion` to make sure that the current Hugo version is greater than the version declared in a theme's config file (`theme.toml`); it also returns the minimum version of Hugo needed as a string (isn't it already a string from the toml metadata?)

None of this except the compare needs more than a string. In the above, there are two strings desired

- full version: "0.38.1-DEV"
- release version: "0.38.1"

And comparing two versions as strings is only moderately harder than comparing them as structs. And if the release engineer
makes an error in declaring the string, it's a bug, and we can have unit tests against that.

All that code could be reduced to this

```go
func CurrentHugoReleaseVersion() string {
    return "0.38.1"
}
func CurrentHugoVersion() string {
    return CurrentHugoReleaseVersion() + "-DEV"
}
func CompareHugoVersions(check, current string) bool {
    // code to compare two strings as versions
}
```

and there would be a few unit tests on the functions to make sure that bad values aren't used;
e.g. that the suffix is appropriate (if it exists), that the version number is only a triple
if the third component is non-zero, and that comparison works.

Yes, there is no "type safety" on the version data. But it's a string. In this particular case,
we just get pain from trying to wrap it.

I haven't made this change yet, but I want to, since it's things like this that make code grow
without bounds.

## Config

Start by creating command-line processors with Cobra.

If we have a command that loads config from config files:

- if `--config` was specified, load all the mentioned config files.
- if `--config` was not specified, look for `config.[toml|yaml|yson]` in working directory.

Merge in flags from the active command-line processor.

Add default config.

Access all config through Viper accessors.

I don't see the point in exactly replicating the existing Hugo code's behavior with config. It
should be general-case, and it's not.

E.g. instead of this

```go
    for _, cmdV := range c.subCmdVs {
        c.initializeFlags(cmdV)
    }
```

where `cmdeer.initializeFlags` cherry-picks among a set of flags,
just do this

```go
    for _, cmdV := range c.subCmdVs {
        config.BindPFlags(cmdV.Flags()) // bind all flags to viper - why not?
        config.BindPFlags(cmdV.PersistentFlags()) // bind all persistent flags to viper - why not?
        c.initializeFlags(cmdV)
    }
```

where we add all flags to the config, under their own names. We still need aliases, for the
case where command-line flags are named differently from config file keys.

### Updating viper

The ability to merge config files was added to Viper, but whole keys are replaced at the base level,
which users find annoying. Maybe there should be an option to merge keys instead of replace? Certainly,
replacing is very predictable, whereas merging has all kinds of edge cases (there's no easy way to
delete sub-keys, for example).

Maybe this is too weird and should be avoided, just like config values replace default values en-masse.

Keys are case-insensitive, but it would be nice to be case-preserving. That would have to be done
with a separate table, because the lookup tables themselves need lower-case keys (there's no way
to override map behavior in Go as there is in other languages).
