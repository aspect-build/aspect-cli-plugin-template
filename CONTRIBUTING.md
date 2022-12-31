## Developing

To try the plugin, first check that you have the most recent [aspect cli release](https://github.com/aspect-build/aspect-cli/releases) installed.

First build the plugin from source:

```bash
% bazel build ...
```

Note that the `.aspect/cli/plugins.yaml` file has a reference to the path under `bazel-bin` where the plugin binary was just written.
On the first build, you'll see a warning printed that the plugin doesn't exist at this path.
This is just the development flow for working on plugins; users will reference the plugin's releases which are downloaded for them automatically.

Now just run `aspect`. You should see that `hello-world` appears in the help output. This shows that our plugin was loaded and contributed a custom command to the CLI.

```
Usage:
  aspect [command]

Custom Commands from Plugins:
  hello-world        Print 'Hello World!' to the command line.
```

## Releasing

Just push a tag to this GitHub repo.
The actions integration will create a release.
