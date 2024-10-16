<p align="center">
<a href="https://github.com/mojoaar/md-app"><img src="https://img.shields.io/github/last-commit/mojoaar/md-app"></a>
<a href="https://github.com/mojoaar/md-app"><img src="https://img.shields.io/github/contributors/mojoaar/md-app"></a>
</p>
<p align="center">
<a href="https://technet.cc"><img src="https://img.shields.io/badge/technet.cc-Blog-blue"></a>
<a href="https://twitter.com/mojoaar"><img src="https://img.shields.io/twitter/follow/mojoaar?style=social"></a>
</p>

# File Creator (md)
File creator is a (simple) command-line utility for creating new markdown files from templates.

## Usage

*General*
```
Markdown File Creator v1.0.0
Author: Morten Johansen (mojoaar)

A tool for creating markdown files and managing templates.

Usage Examples:
  Create a new template:
    md template create my-template

  List all templates:
    md template list

  Create a new note using the default template:
    md note -t "My Note Title"

  Create a new note with a custom name and template:
    md note -t "My Note Title" -n my-custom-note -m my-template

  Create a new note with tags:
    md note -t "My Note Title" -g tag1,tag2,tag3

  List all notes with their tags:
    md list

Usage:
  md [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        List all notes with their tags
  note        Create a new note
  template    Manage templates

Flags:
  -h, --help      help for md
  -v, --version   version for md

Use "md [command] --help" for more information about a command.
```

*Create a new template*
```
md template create my-template
```

*Show all available templates (will also create default.yaml if it is not in the templates directory)*
```
md template list
```

*Create a new note*
```
md note -t "My Note Title"
```

*Version & help information*
```
Use md --version and md --help to get version and help information.
```

## Changelog
* 1.0.0 - 2024-10-08
  * Initial release of app
