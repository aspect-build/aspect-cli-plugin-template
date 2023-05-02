# Template: aspect-cli plugin in Go

This repo provides the fastest way to make a plugin for the [aspect cli](https://aspect.build/cli).

It contains a plugin written in Go, with a GitHub actions CI/CD pipeline to release it.

More details about aspect cli plugins is on the [plugin documentation](https://docs.aspect.build/aspect-build/aspect-cli/5.0.3/docs/help/topics/plugins.html)

## Instructions

Create a new repo with the green "Use this template" button above.
Then in your repo...

1. Find-and-replace `hello-world` with your plugin name.
1. Find-and-replace `github.com/aspect-build/aspect-cli-plugin-template` with the name of your Go module. See <https://go.dev/doc/modules/developing>
1. Start coding on your features!
1. Delete everything above the SNIP line below, and replace it with info about your plugin.
1. Once satisfied, create a tag `vX.Y.Z` and push it to trigger the GitHub Action that will create the release.
1. Navigate to the releases tab of your repository, you'll find a draft release based on your tag. Publish it to complete the release process.
1. Want to share your plugin with other developers? Consider adding it to the plugin catalog, by sending a PR editing the `plugins.json` file located in the public [Aspect CLI plugins registry](https://github.com/aspect-build/aspect-cli/tree/main/docs/plugins).

### Local development

1. Run `bazel build :dev` to build your plugin
1. In an existing codebase using aspect-cli, add the plugin like following:

```
# .aspect/cli/config.yaml
plugins:
  - name: <your-plugin-name>
    from: <relative path to the plugin>/bazel-bin/plugin
```

1. You're all set, just remember to build the `:dev` target whenever you change code in the plugin.

### Working with BEP events 

If you want to work with `BEPEventCallback(event *buildeventstream.BuildEvent) error` but you're not very familiar with the event structures, 
you can run the Bazel command you're looking to augment with the `--build_event_json_file=bep.json` and inspect the resulting `bep.json` to 
get a rough understanding about where you should look into.

See also the [Bazel user guide page about BEP](https://bazel.build/remote/bep).

---------- %<  SNIP %< ------------

# My Plugin

This is a plugin for the Aspect CLI.

## Demo

> TODO: Consider showing off your new plugin with a little animated demo of your terminal! We highly recommend [asciinema](https://asciinema.org).
