<p align="center">
<a href="https://github.com/mojoaar/md-app"><img src="https://img.shields.io/github/last-commit/mojoaar/md-app"></a>
<a href="https://github.com/mojoaar/md-app"><img src="https://img.shields.io/github/contributors/mojoaar/md-app"></a>
</p>
<p align="center">
<a href="https://technet.cc"><img src="https://img.shields.io/badge/technet.cc-Blog-blue"></a>
<a href="https://twitter.com/mojoaar"><img src="https://img.shields.io/twitter/follow/mojoaar?style=social"></a>
</p>

# File Creator (md)
File creator is a command-line utility for creating new markdown files from templates.

## Usage

*General*
```
Create a new template:
  md -type template -name <template_name>
Show all available templates:
  md -type template -show
Create a new note:
  md -type note -title <note_title> [-name <note_name>] [-template <template_name>]
```

*Create a new template*
```
md -type template -name default
```

*Show all available templates (will also create default.yaml if it is not in the templates directory)*
```
md -type template -show
```

*Example response*
```
Available template files:
- default
```

*Create a new note*
```
md -type note -title "My First Note" -template default
```

*Version & help information*  
```
Use md -v and md -h to get version and help information.
```

## Changelog
* 1.0.0 - 2024-10-08
  * Initial release of app