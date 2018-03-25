Notes about Hugo commands.

## Common build flags

Is there a better way to manage Hugo config? Does this much of it need to
be exposed on the command-line as individual options? Maybe much of it
could move to repeated key/value pairs?

`hugo --config=K:V --config=K:V ...`

This would be more explicit; keys would match what's in the default config
or in a config file. And the answer would be to make online help for config
directly available, which could be part of the default config stored in the
code.

## `hugo check`

This has never existed since the very start of the Hugo project. Maybe this
command should be gently retired, as it will probably never have any useful
function. What would this do that wouldn't be covered by Hugo test code?
