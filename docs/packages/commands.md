Notes about Hugo commands.

These are the commands that exist

- hugo
- hugo benchmark
- hugo check
  - hugo check ulimit (Darwin-only)
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
- hugo release (hidden)
- hugo server
- hugo version

I think there should be a `hugo build` and that `hugo` is just an alias to
this command. This both preserves existing behavior and is a better hint
to novice users. The biggest reason to do this is so that typing `hugo`
by itself shows help if it's not clearly a build command - at the moment,
it complains about not finding a config file, which is very new-user-unfriendy.

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
        RunE:  h.hugo,
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
there is only one of these and it is a parameter to every command-line handler function. For
now, these are unbound, and must be fetched with a Get call.
