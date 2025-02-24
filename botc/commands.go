/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package botc

import (
	"github.com/cyb3rko/matrix-botc/util"
	"gopkg.in/yaml.v3"
	"maunium.net/go/mautrix/event"
	"os"
	"regexp"
	"strings"
)

const prefixRegex = "^![a-z0-9]{2,10}$"
const commandRegex = "^[a-z0-9]{1,20}$"

var commandPrefix string
var commandMapping map[string]Command
var selfDisclosureFunction SelfDisclosureFunction
var topLevelHelpFunction CommandHelpFunction

type Config struct {
	Prefix                 string
	Mapping                map[string]Command
	SelfDisclosureFunction SelfDisclosureFunction
	HelpFunction           CommandHelpFunction
}

type Command struct {
	EndCommand        bool
	RequiredArguments int
	ProcessFunction   CommandFunction
	HelpFunction      CommandHelpFunction
	Subcommands       map[string]Command
}

type SelfDisclosure struct {
	Name        string
	Author      string
	Version     string
	About       string
	Interactive bool
}

type SelfDisclosureFunction func(disclosure *SelfDisclosure)

type CommandFunction func(evt *event.Event, args []string)

type CommandHelpFunction func()

// RegisterCommands registers the commands which should be processed by the botc library
//
//goland:noinspection GoUnusedExportedFunction
func RegisterCommands(config *Config) {
	print("Botc: Registering commands for this bot")
	commandPrefix = validatePrefix(config.Prefix)
	validateDisclosureConfig(config.SelfDisclosureFunction)
	cmdRegex := regexp.MustCompile(commandRegex)
	valid, err := checkCommands(cmdRegex, config.Mapping, 0)
	if !valid {
		panic(err)
	}
	commandMapping = config.Mapping
	selfDisclosureFunction = config.SelfDisclosureFunction
	topLevelHelpFunction = config.HelpFunction
	print("Botc: Config check complete, commands are registered")
}

func validatePrefix(prefix string) string {
	prefix = strings.ToLower(prefix)
	if !regexp.MustCompile(prefixRegex).MatchString(prefix) {
		panic(format("Botc: Prefix '%s' does not match allowed format '%s'", prefix, prefixRegex))
	}
	return prefix
}

func validateDisclosureConfig(selfDisclosureFunction SelfDisclosureFunction) {
	if selfDisclosureFunction == nil {
		return
	}
	config, err := os.ReadFile("disclosure.yaml")
	if err != nil {
		panic(format("Botc: Error reading disclosure.yaml: %s", err))
	}
	data := make(map[interface{}]interface{})
	err = yaml.Unmarshal(config, &data)
	if err != nil {
		panic(format("Botc: Error processing disclosure.yaml: %s", err))
	}
	valid, missingKey := util.HasMapKeys(data, []string{"name", "author", "version", "about", "interactive"})
	if !valid {
		panic(format("Botc: disclosure.yml does not contain key '%s'", missingKey))
	}
}

// ProcessCommandChain triggers the processing on user input and executes the registered command functions
//
// It returns whether the input was valid or not
//
//goland:noinspection GoUnusedExportedFunction
func ProcessCommandChain(input string, evt *event.Event) bool {
	input = strings.TrimSpace(input)
	lowerInput := strings.ToLower(input)
	if !strings.HasPrefix(lowerInput, commandPrefix+" ") && lowerInput != commandPrefix {
		// not related to registered prefix
		if lowerInput == "!bots" {
			// self-disclosure trigger
			processSelfDisclosure()
			return true
		}
		return false
	}
	commandParts := strings.Split(input, " ")[1:]
	if len(commandParts) == 0 {
		print("Botc: Received only prefix; showing top-level help page")
		topLevelHelpFunction()
		return false
	}
	_, valid := commandMapping[strings.ToLower(commandParts[0])]
	if !valid {
		printf("Botc: Received unregistered command '%s'; showing top-level help page", commandParts[0])
		topLevelHelpFunction()
		return false
	}
	printf("Botc: Received command '%s': '%s'", commandParts[0], strings.Join(commandParts[1:], " "))
	return processCommand(commandMapping, commandParts, topLevelHelpFunction, evt)
}

// Process command and subcommands recursisvely
func processCommand(
	mapping map[string]Command,
	commandParts []string,
	parentHelpFunction CommandHelpFunction,
	evt *event.Event,
) bool {
	command, valid := mapping[strings.ToLower(commandParts[0])]
	if !valid {
		// command / subcommand not found
		printf("Botc: (Sub)Command '%s' not found; showing parent help page", strings.ToLower(commandParts[0]))
		parentHelpFunction()
		return false
	}
	subcommandParts := commandParts[1:]

	if !command.EndCommand {
		// intermediate command, not the last part of the command chain
		if len(subcommandParts) >= command.RequiredArguments {
			return processCommand(command.Subcommands, subcommandParts, command.HelpFunction, evt)
		} else {
			// not enough subcommands & arguments
			printf("Botc: Command chain '%[1]s' contains not enough subcommands & arguments; showing '%[1]s' help page", strings.Join(commandParts, " "))
			command.HelpFunction()
			return true
		}
	} else {
		// last part of the command chain
		if len(subcommandParts) == command.RequiredArguments {
			command.ProcessFunction(evt, subcommandParts)
			return true
		} else {
			// unexpected amount of arguments
			printf("Botc: Command '%s' contains not enough arguments", strings.Join(commandParts, " "))
			parentHelpFunction()
			return false
		}
	}
}

func processSelfDisclosure() {
	if selfDisclosureFunction != nil {
		config, err := os.ReadFile("disclosure.yaml")
		if err != nil {
			return
		}
		data := make(map[interface{}]interface{})
		err = yaml.Unmarshal(config, &data)
		if err != nil {
			return
		}
		valid, _ := util.HasMapKeys(data, []string{"name", "author", "version", "about", "interactive"})
		if !valid {
			return
		}
		selfDisclosure := SelfDisclosure{
			Name:        data["name"].(string),
			Author:      data["author"].(string),
			Version:     data["version"].(string),
			About:       data["about"].(string),
			Interactive: data["interactive"].(bool),
		}
		selfDisclosureFunction(&selfDisclosure)
	}
}

// Check recursively for valid commands and subcommands; maximum depth of 10 allowed
func checkCommands(cmdRegex *regexp.Regexp, mapping map[string]Command, depth int) (bool, string) {
	if depth == 10 {
		return false, "Botc: Too many subcommands registered, maximum of 10 allowed"
	}
	for name, command := range mapping {
		name = strings.ToLower(name)
		if !cmdRegex.MatchString(name) {
			return false, format("Botc: Command '%s' does not match allowed format '%s'", name, cmdRegex)
		}
		if valid, err := checkCommand(name, command); !valid {
			return false, err
		}
		if command.Subcommands != nil {
			if valid, err := checkCommands(cmdRegex, command.Subcommands, depth+1); !valid {
				return false, err
			}
		}
	}
	return true, ""
}

// Check a command for misconfiguration
func checkCommand(name string, command Command) (bool, string) {
	if !command.EndCommand {
		// intermediate command
		if command.RequiredArguments < 1 {
			return false, format("Botc: Intermediate command '%s' must have at least one argument", name)
		}
		if command.ProcessFunction != nil {
			return false, format("Botc: Intermediate command '%s' can not have a process function", name)
		}
		if command.HelpFunction == nil {
			return false, format("Botc: Intermediate command '%s' must have a help function", name)
		}
		if command.Subcommands == nil {
			return false, format("Botc: Intermediate command '%s' must have subcommands", name)
		}
	} else {
		// end command
		if command.RequiredArguments < 0 {
			return false, format("Botc: End command '%s' has negative amount of arguments? What does that mean? Use 0 instead", name)
		}
		if command.ProcessFunction == nil {
			return false, format("Botc: End command '%s' must have a process function", name)
		}
		if command.HelpFunction != nil {
			return false, format("Botc: End command '%s' can not have a help function", name)
		}
		if command.Subcommands != nil {
			return false, format("Botc: End command '%s' can not have subcommands", name)
		}
	}
	return true, ""
}
