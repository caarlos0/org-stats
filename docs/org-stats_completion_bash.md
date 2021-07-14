# org-stats completion bash

generate the autocompletion script for bash

## Synopsis


Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(org-stats completion bash)

To load completions for every new session, execute once:

### Linux:

	org-stats completion bash > /etc/bash_completion.d/org-stats

### macOS:

	org-stats completion bash > /usr/local/etc/bash_completion.d/org-stats

You will need to start a new shell for this setup to take effect.
  

```
org-stats completion bash
```

## Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

## See also

* [org-stats completion](org-stats_completion.md)	 - generate the autocompletion script for the specified shell

