# Template for aspect-cli plugins

This repo provides the fastest way to make a plugin for the [aspect cli].

It contains a plugin written in Go, with a GitHub actions CI/CD pipeline to release it.

More details about aspect cli plugins is on the [plugin documentation].

## Instructions

Create a new repo with the green "Use this template" button above.
Then in your repo...

1. Replace `aspect_plugin_hello_world` with ...
1. Replace `hello_world` with ...
1. Replace `github.com/aspect-build/aspect-cli-plugin-template` with the name of your Go module. See <https://go.dev/doc/modules/developing>
1. Delete everything above the SNIP line below, and start coding on your features!

---------- %<  SNIP %< ------------

# aspect_plugin_hello_world

This is a plugin for the Aspect CLI.

## Developing

To try the plugin, first check that you have the most recent [aspect cli release] installed.

First build it. We assume you have installed [bazelisk] on your $PATH as `bazel`.

```bash
% bazel build ...
```

> Note that the `.aspectplugins` file has a reference to the path under `bazel-bin` where the plugin binary was just written.

Now just run `aspect`. You should see:

```
Usage:
  aspect [command]

Available Commands:
  ...
  hello-world   Print 'Hello World!' to the command line.
```

This shows that our plugin was loaded and contributed a custom command to the CLI.

## Releasing

Just push a tag to your GitHub repo.
The actions integration will create a release.

[bazelisk]: https://bazel.build/install/bazelisk
[aspect cli]: https://aspect.build
[plugin documentation]: https://aspect.build/help/topics/plugins
[aspect cli release]: https://github.com/aspect-build/aspect-cli/releases