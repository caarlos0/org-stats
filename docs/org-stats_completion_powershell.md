# org-stats completion powershell

generate the autocompletion script for powershell

## Synopsis


Generate the autocompletion script for powershell.

To load completions in your current shell session:

	org-stats completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
org-stats completion powershell [flags]
```

## Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

## See also

* [org-stats completion](org-stats_completion.md)	 - generate the autocompletion script for the specified shell

