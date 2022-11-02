package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"gopkg.in/yaml.v3"
)

type ApplicationConfig struct {
	Name              string
	Path              string
	KeybindingPattern string
}

type Config struct {
	Applications []ApplicationConfig
}

type Keybinding struct {
	Binding    string
	Definition string
}

func main() {
	err := ensureConfig()
	if err != nil {
		log.Fatalf("Configuration error: %s", err.Error())
	}

	configDir, err := configurationDirectory()
	if err != nil {
		log.Fatalf("Error acquiring configuration directory: %s", err.Error())
	}

	configPath := configurationPath(configDir)

	confBytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err.Error())
	}

	var config Config
	err = yaml.Unmarshal(confBytes, &config)
	if err != nil {
		log.Fatalf("Error processing config file: %s", err.Error())
	}

	var (
		choice int
		argNum = len(os.Args)
	)

	if argNum == 1 {
		choice, err = getChoiceFromTerminal(&config)
	} else if argNum == 2 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			fmt.Printf(
				"\n%s [application name]\n\n"+
					"Lists the keybindings of the application name given as an argument.\n"+
					"Configs can be found at: %s\n\n",
				os.Args[0], configPath,
			)
			os.Exit(0)
		}

		choice, err = getChoiceFromCommandLine(&config)
	} else {
		err = fmt.Errorf("Zero or one command line expected, got %d", argNum)
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error retrieving the home directory")
	}

	configPath = filepath.Join(userHomeDir, config.Applications[choice].Path)
	confBytes, err = os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading config file (%s): %s", configPath, err.Error())
	}

	keybindings := []Keybinding{}
	pattern := config.Applications[choice].KeybindingPattern
	scanner := bufio.NewScanner(strings.NewReader(string(confBytes)))

	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatalf("Regexp could not be compiled (%s): %s", pattern, err.Error())
	}

	for scanner.Scan() {
		line := scanner.Text()

		matched := re.FindAllSubmatch([]byte(line), -1)

		if len(matched) > 0 {
			for _, match := range matched {
				keybindings = append(
					keybindings,
					Keybinding{
						Binding:    string(match[1]),
						Definition: string(match[2]),
					},
				)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error parsing config file (%s): %s", configPath, err.Error())
	}

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Binding", "Definition")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, binding := range keybindings {
		tbl.AddRow(binding.Binding, binding.Definition)
	}

	fmt.Println("")
	tbl.Print()
	fmt.Println("")
}

func defaultConfig() Config {
	return Config{
		Applications: []ApplicationConfig{
			{
				Name:              "i3",
				Path:              ".config/i3/config",
				KeybindingPattern: "bindsym ([a-zA-Z0-9$+]+) (.*)",
			},
			{
				Name:              "vim",
				Path:              ".vimrc",
				KeybindingPattern: "(?:map|nmap|nnoremap|tnoremap) ((?:[a-zA-Z0-9<>-]|\\\\p{Punct})+) (.*)",
			},
			{
				Name:              "vifm",
				Path:              ".config/vifm/vifmrc",
				KeybindingPattern: "nnoremap ([a-zA-Z0-9<>,]+) (.*)",
			},
		},
	}
}

func ensureConfig() error {
	configDir, err := configurationDirectory()

	if _, err := os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(configDir, os.ModePerm); err != nil {
			return fmt.Errorf("Couldn't create config directory (%s): %w", configDir, err)
		}
	}

	var configPath = configurationPath(configDir)

	if _, err = os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		var config = defaultConfig()
		bytes, err := yaml.Marshal(&config)
		if err != nil {
			return fmt.Errorf("Error marshalling default config: %w", err)
		}

		err = os.WriteFile(configPath, bytes, 0644)
		if err != nil {
			return fmt.Errorf("Error creating default config: %w", err)
		}
	}

	return nil
}

func configurationDirectory() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Couldn't retrieve user directory: %w", err)
	}

	return filepath.Join(userHomeDir, ".config", "keybindings"), nil
}

func configurationPath(configDir string) string {
	return filepath.Join(configDir, "config.yml")
}

func getChoiceFromTerminal(config *Config) (choice int, err error) {
	fmt.Printf("\nChoose application:\n\n")
	for i, app := range config.Applications {
		fmt.Printf("(%d)  %s\n", i, app.Name)
	}

	var inputOK = false

	for !inputOK {
		fmt.Printf("\nChoice: ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return 0, fmt.Errorf("An error occured while reading input: %s", err.Error())
		}

		choice, err = strconv.Atoi(strings.TrimSuffix(input, "\n"))
		if err != nil {
			fmt.Printf("\nInvalid input: %s\n", err.Error())
		} else {
			inputOK = true
		}

		if inputOK && (choice < 0 || choice >= len(config.Applications)) {
			fmt.Printf("Choice should be between 0 and %d\n", len(config.Applications)-1)
			inputOK = false
		}
	}

	return
}

func getChoiceFromCommandLine(config *Config) (choice int, err error) {
	appName := os.Args[1]

	for i, appConf := range config.Applications {
		if appConf.Name == appName {
			return i, nil
		}
	}

	return 0, fmt.Errorf("Unrecognized application name: %s", appName)
}
