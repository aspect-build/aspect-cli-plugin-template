/*
Copyright © 2021 Aspect Build Systems Inc

Not licensed for re-use.
*/

package main

import (
	goplugin "github.com/hashicorp/go-plugin"

	buildeventstream "aspect.build/cli/bazel/buildeventstream/proto"
	"aspect.build/cli/pkg/ioutils"
	"aspect.build/cli/pkg/plugin/sdk/v1alpha2/config"
)

func main() {
	// config.NewConfigFor accepts a plugin implementation and returns the go-plugin
	// configuration required to serve the plugin to the CLI core.
	goplugin.Serve(config.NewConfigFor(NewDefaultPlugin()))
}

// YourPlugin implements an aspect CLI plugin.
type YourPlugin struct{}

// NewDefaultPlugin creates a new YourPlugin with the default
// dependencies.
func NewDefaultPlugin() *YourPlugin {
	return NewPlugin()
}

// NewPlugin creates a new YourPlugin, allowing dependencies to be
// injected.
func NewPlugin() *YourPlugin {
	return &YourPlugin{}
}

func (plugin *YourPlugin) BEPEventCallback(event *buildeventstream.BuildEvent) error {
	// Process build events here
	return nil
}

func (plugin *YourPlugin) PostBuildHook(
	isInteractiveMode bool,
	promptRunner ioutils.PromptRunner,
) error {
	// Perform a post build action
	return nil
}

func (plugin *YourPlugin) PostTestHook(
	isInteractiveMode bool,
	promptRunner ioutils.PromptRunner,
) error {
	// Perform a post test action
	return nil
}

func (plugin *YourPlugin) PostRunHook(
	isInteractiveMode bool,
	promptRunner ioutils.PromptRunner,
) error {
	// Perform a post run action
	return nil
}
