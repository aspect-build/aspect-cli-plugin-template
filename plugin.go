package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/manifoldco/promptui"

	"github.com/aspect-build/aspect-cli/bazel/buildeventstream"
	"github.com/aspect-build/aspect-cli/bazel/command_line"
	"github.com/aspect-build/aspect-cli/pkg/ioutils"
	"github.com/aspect-build/aspect-cli/pkg/plugin/sdk/v1alpha4/config"
	aspectplugin "github.com/aspect-build/aspect-cli/pkg/plugin/sdk/v1alpha4/plugin"
)

// main starts up the plugin as a child process of the CLI and connects the gRPC communication.
func main() {
	goplugin.Serve(config.NewConfigFor(&HelloWorldPlugin{
		besChan: make(chan orderedBuildEvent, 100),
	}))
}

// HelloWorldPlugin declares the fields on an instance of the plugin.
type HelloWorldPlugin struct {
	// Base gives default implementations of the plugin methods, so implementing them below is optional.
	// See the definition of aspectplugin.Base for more methods that can be implemented by the plugin.
	aspectplugin.Base
	// This plugin will store some state from the Build Events for use at the end of the build.
	command_line.CommandLine

	besOnce             sync.Once
	besChan             chan orderedBuildEvent
	besHandlerWaitGroup sync.WaitGroup
}

type orderedBuildEvent struct {
	event          *buildeventstream.BuildEvent
	sequenceNumber int64
}

// CustomCommands contributes a new 'hello-world' command alongside the built-in ones like 'build' and 'test'.
func (plugin *HelloWorldPlugin) CustomCommands() ([]*aspectplugin.Command, error) {
	return []*aspectplugin.Command{
		aspectplugin.NewCommand(
			"hello-world",
			"Print 'Hello World!' to the command line.",
			"Print 'Hello World!' to the command line. Echo any given argument.",
			func(ctx context.Context, args []string, bazelStartupArgs []string) error {
				fmt.Println("Hello World!")
				fmt.Print("Arguments passed to command: ")
				fmt.Println(args)
				return nil
			},
		),
	}, nil
}

// BEPEventCallback subscribes to all Build Events, and lets our logic react to ones we care about.
func (plugin *HelloWorldPlugin) BEPEventCallback(event *buildeventstream.BuildEvent, sequenceNumber int64) error {
	plugin.besChan <- orderedBuildEvent{event: event, sequenceNumber: sequenceNumber}

	plugin.besOnce.Do(func() {
		plugin.besHandlerWaitGroup.Add(1)
		go func() {
			defer plugin.besHandlerWaitGroup.Done()
			var nextSn int64 = 1
			eventBuf := make(map[int64]*buildeventstream.BuildEvent)
			for o := range plugin.besChan {
				if o.sequenceNumber == 0 {
					// Zero is an invalid squence number. Process the event since we can't order it.
					if err := plugin.BEPEventHandler(o.event); err != nil {
						log.Printf("error handling build event: %v\n", err)
					}
					continue
				}

				// Check for duplicate sequence numbers
				if _, exists := eventBuf[o.sequenceNumber]; exists {
					log.Printf("duplicate sequence number %v\n", o.sequenceNumber)
					continue
				}

				// Add the event to the buffer
				eventBuf[o.sequenceNumber] = o.event

				// Process events in order
				for {
					if orderedEvent, exists := eventBuf[nextSn]; exists {
						if err := plugin.BEPEventHandler(orderedEvent); err != nil {
							log.Printf("error handling build event: %v\n", err)
						}
						delete(eventBuf, nextSn) // Remove processed event
						nextSn++                 // Move to the next expected sequence
					} else {
						break
					}
				}
			}
		}()
	})

	return nil
}

func (plugin *HelloWorldPlugin) BEPEventHandler(event *buildeventstream.BuildEvent) error {
	switch event.Payload.(type) {
	case *buildeventstream.BuildEvent_StructuredCommandLine:
		commandLine := *event.GetStructuredCommandLine()
		if commandLine.CommandLineLabel == "canonical" {
			plugin.CommandLine = commandLine
		}
	}
	return nil
}

// PostBuildHook will be called at the end of an `aspect build` execution, after Bazel completes.
func (plugin *HelloWorldPlugin) PostBuildHook(
	isInteractiveMode bool,
	promptRunner ioutils.PromptRunner,
) error {
	// Close the build events channel
	close(plugin.besChan)

	// Wait for all build events to come in
	if !waitGroupWithTimeout(&plugin.besHandlerWaitGroup, 60*time.Second) {
		log.Printf("timed out waiting for BES events\n")
	}

	// We condition prompting on whether there's an interactive user to engage with.
	if isInteractiveMode {
		// The manifoldco/promptui library creates many styles of interactive prompts.
		// Check out the examples: https://github.com/manifoldco/promptui/tree/master/_examples
		prompt := promptui.Prompt{
			Label:     "Thanks for trying the hello-world plugin! Would you like to see the command that was run",
			IsConfirm: true,
		}
		// Since the prompt is a boolean, any non-nil error should represent a NO.
		if _, err := promptRunner.Run(prompt); err == nil {
			plugin.printTargetPattern()
		}
	}
	return nil
}

// PostTestHook satisfies the Plugin interface. In this case, it just calls the PostBuildHook.
func (plugin *HelloWorldPlugin) PostTestHook(
	isInteractiveMode bool,
	promptRunner ioutils.PromptRunner,
) error {
	return plugin.PostBuildHook(isInteractiveMode, promptRunner)
}

// PostRunHook satisfies the Plugin interface. In this case, it just calls the PostBuildHook.
func (plugin *HelloWorldPlugin) PostRunHook(
	isInteractiveMode bool,
	promptRunner ioutils.PromptRunner,
) error {
	return plugin.PostBuildHook(isInteractiveMode, promptRunner)
}

// printTargetPattern is just representative of some logic a plugin might want to perform on the data collected.
func (plugin *HelloWorldPlugin) printTargetPattern() {
	for _, section := range plugin.CommandLine.Sections {
		fmt.Fprintf(os.Stdout, "%s\n", section.SectionLabel)
		if section.SectionLabel == "residual" {
			switch f := section.SectionType.(type) {
			case *command_line.CommandLineSection_ChunkList:
				fmt.Fprintf(os.Stdout, "target pattern was %s\n", f.ChunkList.Chunk[0])
			}
		}
	}
}

// waitGroupWithTimeout waits for a WaitGroup with a specified timeout.
func waitGroupWithTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	done := make(chan struct{})

	// Run a goroutine to close the channel when WaitGroup is done
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// WaitGroup finished within timeout
		return true
	case <-time.After(timeout):
		// Timeout occurred
		return false
	}
}
