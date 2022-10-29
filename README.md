# Keybindings

This application parses application configs for keyboard shortcuts and echos them.

It uses information given in its config file for each application:

- the application name: it is used in the command line arguments to tell Shortcuts
	which application's config to look for
- the config's path: the path to the config file (either absolute path
	or relative to the user's home directory
- the syntax: the regex that matches the lines that define shortcuts
	- the regex must contain two unnamed capturing groups: the first is the keyboard
		shortcut, the second is the command that's executed
	- the regex can contain any number of additional non-capturing groups
		- Form: `(?:regex)`

The config file is searched in: `$HOME/.config/shortcuts/config.yml`.
If it is not found there then the directory is created and a default config
file is placed there.

The default configuration:

```yaml
applications:
    - name: i3
      path: .config/i3/config
      keybindingpattern: bindsym ([a-zA-Z0-9$+]+) (.*)
    - name: vim
      path: .vimrc
      keybindingpattern: (?:map|nmap|nnoremap|tnoremap) ((?:[a-zA-Z0-9<>]|\\p{Punct})+) (.*)
    - name: vifm
      path: .config/vifm/vifmrc
      keybindingpattern: nnoremap ([a-zA-Z0-9<>,]+) (.*)
```

## Compilation

The application needs to be compiled. For this, a working Go build environment
is needed. Please refer to your distribution's package manager or visit
[this][1] page for information on how to install Go.

>	**Note**
>
>	Go can typically be installed using package managers, although it might not
>	be the latest version that's available (especially if the package manager
>	serves a point release distribution). Some install commands for popular
>	package managers are:
>
>	Arch/Manjaro: `packman -S go`
>	Fedora: `dnf install golang`
>	Ubuntu: `apt-get install golang`

Once Go is available on the system, the compilation can be performed by
issuing:

```bash
go build -mod=mod -o build/keybindings ./cmd/keybindings
```

or by using the attached Makefile:

```bash
make all
```

The binary will be available in the `build` folder of the repositry root.

## Installation

The application can also be installed from `github` using the following
command:

```bash
go install github.com/nagygr/keybindings/cmd/keybindings@latest
```

>	**Note**
>
>	Please note, that the command above also requires Go to be installed on the
>	system. Please see details about Go installation above.

[1]: https://go.dev/doc/install

