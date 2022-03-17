package main

import (
	"context"
	"fmt"
	"os"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/manifoldco/promptui"

	"aspect.build/cli/bazel/buildeventstream"
	"aspect.build/cli/bazel/command_line"
	"aspect.build/cli/pkg/bazel"
	"aspect.build/cli/pkg/ioutils"
	"aspect.build/cli/pkg/plugin/sdk/v1alpha3/config"
	aspectplugin "aspect.build/cli/pkg/plugin/sdk/v1alpha3/plugin"
)

func main() {
	goplugin.Serve(config.NewConfigFor(NewDefaultPlugin()))
}

type HelloWorldPlugin struct {
	aspectplugin.Base
	command_line.CommandLine
}

func NewDefaultPlugin() *HelloWorldPlugin {
	return NewPlugin()
}

func NewPlugin() *HelloWorldPlugin {
	return &HelloWorldPlugin{}
}

func (plugin *HelloWorldPlugin) CustomCommands() ([]*aspectplugin.Command, error) {
	return []*aspectplugin.Command{
		aspectplugin.NewCommand(
			"hello-world",
			"Print 'Hello World!' to the command line.",
			"Print 'Hello World!' to the command line. Echo any given argument. Then run a 'bazel help'",
			func(ctx context.Context, args []string, bzl bazel.Bazel) error {
				fmt.Println("Hello World!")
				fmt.Print("Arguments passed to command: ")
				fmt.Println(args)
				fmt.Println("Going to run: 'bazel help'")

				bzl.Spawn([]string{"help"})

				return nil
			},
		),
	}, nil
}
// BEPEventCallback satisfies the Plugin interface. It process all the analysis
// failures that represent a visibility issue, collecting them for later
// processing in the post-build hook execution.
func (plugin *HelloWorldPlugin) BEPEventCallback(event *buildeventstream.BuildEvent) error {
	switch event.Payload.(type) {
		case *buildeventstream.BuildEvent_StructuredCommandLine:
			plugin.CommandLine = *event.GetStructuredCommandLine()
	}
	return nil
}

func (plugin *HelloWorldPlugin) PostBuildHook(
	isInteractiveMode bool,
	promptRunner ioutils.PromptRunner,
) error {
	if isInteractiveMode {
		prompt := promptui.Prompt{
			Label:     "Would you like to see the command that was run",
			IsConfirm: true,
		}
		// Since the prompt is a boolean, any non-nil error should represent a NO.
		if _, err := promptRunner.Run(prompt); err == nil {
			plugin.printTargetPattern()
		}
	}
	return nil
}

func (plugin *HelloWorldPlugin) printTargetPattern() {
	for _, section := range plugin.CommandLine.Sections {
		fmt.Fprintf(os.Stdout, "%s\n", section.SectionLabel)
		if section.SectionLabel == "residual" {
			switch f := section.SectionType.(type) {
			case *command_line.CommandLineSection_ChunkList:
				fmt.Fprintf(os.Stdout, "pattern was %s\n", f.ChunkList.Chunk[0])
			}
		}
	}
}