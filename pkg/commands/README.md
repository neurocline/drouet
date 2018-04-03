Notes about Hugo commands.

These are the commands that exist

- hugo
- hugo benchmark
- hugo check
- hugo config
- hugo convert
  - hugo convert toJSON
  - hugo convert toTOML
  - hugo convert toYAML
- hugo env
- hugo gen
  - hugo gen autocomplete
  - hugo gen doc
  - hugo gen man
  - hugo gen docshelper (hidden)
  - hugo gen chromastyles
- hugo import
  - hugo import jekyll
- hugo list
  - hugo list drafts
  - hugo list expired
  - hugo list future
- hugo new
  - hugo new site
  - hugo new theme
- hugo server
- hugo version

I think there should be a `hugo build` and that `hugo` is just an alias to
this command. This both preserves existing behavior and is a better hint
to novice users. The biggest reason to do this is so that typing `hugo`
by itself shows help if it's not clearly a build command - at the moment,
it complains about not finding a config file, which is very new-user-unfriendy.

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

## `hugo gen`

If the gen code gets large, it should be split up into a new package
`hugo/pkg/commands/gen/.`, leaving all the command-line setup and processing
in `hugp/pkg/commands/gen.go`

## Accessing Hugo state from Hugo commands

The Cobra library doesn't really have provision for passing state in to
command handlers (specified with variants of `Run` in a `cobra.Command`
struct). There are two Go ways to handle this

1. Use method values for the command handlers
2. Use local functions to capture state in the command handlers.

Either way, the state will be held onto until the end of the program, unless
the command handler structures are disposed of, somehow. Since the same code
controls the command handlers and execution, this is reasonable and safe.

Option 1 is a bit tricky in that you can't just add new methods to a type
you don't own. Struct embedding bypasses this.

Option 2 is simpler but requires every command handler to have a local
literal function that wraps the state and calls the real handler (or is
the entire command handler).

## Pattern for Hugo commands

There is a top-level object `Hugo` that holds the common state for all
Hugo commands (logging, config, sites).

Each command has its own command object and command builder. The pattern
looks like this.

```go

// Build Hugo root command.
func buildHugoCommand(hugo *core.Hugo) *hugoCmd {
    h := &hugoCmd{Hugo: hugo}

    h.cmd = &cobra.Command{
        Use:   "hugo",
        Short: "hugo builds your site",
        RunE: h.hugo,
    }

    // Add flags for the "hugo" command
    h.cmd.Flags().BoolVar(&h.renderToMemory, "renderToMemory", false, "render to memory")
    h.cmd.Flags().BoolVarP(&h.watch, "watch", "w", false, "watch filesystem)

    // Add flags shared by builders: "hugo", "hugo server", "hugo benchmark"
    addHugoBuilderFlags(h.cmd)

    return h
}

type hugoCmd struct {
    *core.Hugo
    cmd *cobra.Command

    renderToMemory bool
    watch bool
}

func (h *hugoCmd) hugo(cmd *cobra.Command, args []string) error {
    return nil
}
```

Each command handler function is a method value bound to `RunE`; this means its receiver is
bound to the function set as `RunE`.

This avoids all global variables, at least those required for passing information between the
command-line and the various command-line handler functions. There are a few down sides at the
moment:

- persistent and shared flags aren't mirrored to struct variables
- mildly repetitive code to set up each command-line object

For top-level persistent variables, we could mirror these to the top-level `*core.Hugo` object;
there is only one of these and it is a parameter to every command-line handler function.
