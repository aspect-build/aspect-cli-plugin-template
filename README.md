# Aspect CLI Plugin Template

### Development Example

As an example we will print `"Hello world"` after the bazel build has completed using `your-plugin`.

1. In your IDE, navigate to `plugins/your-plugin/plugin.go`. In this file you will find the hooks that are available to the aspect cli.
One of those hooks will be `PostBuildHook`. This is the hook that we will use in this example. 

2. Add `"fmt"` to the golang imports and within `PostBuildHook` add the following line:
```
fmt.Println("Hello World")
```

3. We now need to build the plugin so that the aspect cli can use it. In your terminal, the following command should work:
```
aspect build //plugins/your-plugin
```

4. We can now focus on the examples folder. Take note of `examples/your-plugin/.aspectplugins` where you will see a `from` field that references the plugin we just built. To test our plugin we will need to navigate to the examples folder in our terminal. If we run a build we should see our `"Hello World"` appear after the build has completed. You can use the following commands:
```
cd examples/your-plugin
aspect build //...
```