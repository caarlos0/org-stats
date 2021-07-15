# org-stats

Get the contributor stats summary from all repos of any given organization

## Synopsis

org-stats can be used to get an overall sense of your org's contributors.

It uses the GitHub API to grab the repositories in the given organization.
Then, iterating one by one, it gets statistics of lines added, removed and number of commits of contributors.
After that, if opted in, it does several searches to get the number of pull requests reviewed by each of the previously find contributors.
Finally, it prints a rank by each category.


Important notes:
* The `--since` filter does not work "that well" because GitHub summarizes thedata by week, so the data is not as granular as it should be.
* The `--include-reviews` only grabs reviews from users that had contributions on the previous step.
* In the `--blacklist` option, 'foo' blacklists both the 'foo' user and 'foo' repo, while 'user:foo' blacklists only the user and 'repo:foo' only the repository.
* The `--since` option accepts all the regular time.Durations Go accepts, plus a few more: 1y (365d), 1mo (30d), 1w (7d) and 1d (24h).

```
org-stats [flags]
```

## Options

```
  -b, --blacklist strings   blacklist repos and/or users
      --csv-path string     path to write a csv file with all data collected
      --github-url string   custom github base url (if using github enterprise)
  -h, --help                help for org-stats
      --include-reviews     include pull request reviews in the stats
  -o, --org string          github organization to scan
      --since string        time to look back to gather info (0s means everything) (default "0s")
      --token string        github api token (default $GITHUB_TOKEN)
      --top int             how many users to show (default 3)
```

## See also

* [org-stats completion](org-stats_completion.md)	 - generate the autocompletion script for the specified shell
* [org-stats version](org-stats_version.md)	 - Prints org-stats version

