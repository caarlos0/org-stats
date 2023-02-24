# org-stats

[![Release](https://img.shields.io/github/release/caarlos0/org-stats.svg?style=for-the-badge)](https://github.com/caarlos0/org-stats/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](/LICENSE.md)
[![Build status](https://img.shields.io/github/actions/workflow/status/caarlos0/org-stats/build.yml?style=for-the-badge)](https://github.com/caarlos0/org-stats/actions?workflow=build)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](http://godoc.org/github.com/caarlos0/org-stats)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)

Get the contributor stats summary from all repos of any given organization

## screenshots

<p align="center">
  <img width="49%" alt="image" src="https://user-images.githubusercontent.com/245435/125717673-163da857-3456-4b98-ab66-29f5bc61e7cf.png">
  <img width="49%" alt="image" src="https://user-images.githubusercontent.com/245435/125717683-821e13cc-3b2f-4c4d-9032-b69eb26bf5c6.png">
</p>


## usage

Check the [docs folder](/docs/org-stats.md).

## install

### macOS

```sh
brew install caarlos0/tap/org-stats
```

### linux

#### snap

```sh
snap install org-stats
```

#### apt

```sh
echo 'deb [trusted=yes] https://apt.fury.io/caarlos0/ /' | sudo tee /etc/apt/sources.list.d/caarlos0.list
sudo apt update
sudo apt install org-stats
```

#### yum

```sh
echo '[caarlos0]
name=caarlos0
baseurl=https://yum.fury.io/caarlos0/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/caarlos0.repo
sudo yum install org-stats
```

#### arch linux

```sh
yay -S org-stats-bin
```

## stargazers over time

[![Stargazers over time](https://starchart.cc/caarlos0/org-stats.svg)](https://starchart.cc/caarlos0/org-stats)

