<!-- [![Build Status](https://travis-ci.org/n8henrie/fixbashhistory.svg?branch=master)](https://travis-ci.org/n8henrie/fixbashhistory) -->
# fixbashhistory

Remove duplicates in bash history file

- Free software: MIT
<!-- - Documentation: https://fixbashhistory.readthedocs.org -->

## Features

- Deduplicates the `.bash_history` file.
- Preserves the timestamp for the *most recent* invocation of a command.
- Preserves the order of (unique) commands with duplicate timestamps.

## Introduction

Even with `HISTCONTROL=ignoreboth:erasedups` set, duplicate commands that are
**not** sequential are preserved. `tail -f ~/.bash_history` in one window, and
watch it while you run in another:

```shell_session
$ echo foo && history -a
foo
$ echo foo && history -a
foo
$ echo bar && history -a
bar
$ echo foo && history -a
foo
```

`awk` is one way to get an idea of how many duplicate commands are in your
`~/.bash_history`.

```shell_session
$ awk '!/^#[[:digit:]]+$/ && seen[$0]++ \
    { count++ } END { print count }' ~/.bash_history/bash_history
83118
```

I obviously have a lot of duplicates, in part from having too many tmux
sessions open and goofing around with `PROMPT_COMMAND` and `history -a`.

My idea in writing this little script was twofold:

1. Learn a little more Go. Which is why you probably shouldn't use this -- I
   don't know what I'm doing.
1. Deduplicate my large `.bash_history` file.

Before you use this tool, I highly recommend that you back up your
`~/.bash_history` file, and keep that backup around for a while in case
something got messed up that isn't immediately obvious.

## Dependencies

- Go >= 1.8

## Quickstart

```shell_session
$ # Consider `history -a` first
$ go get -u https://github.com/n8henrie/fixbashhistory
$ fixbashhistory -history-file ~/.bash_history -outfile bash_history.new
$ wc -l ~/.bash_history bash_history.new
$ ls -lh ~/.bash_history bash_history.new
```

Check out `bash_history.new`, and if you like what you see:

```shell_session
$ cp ~/.bash_history{,.bak}
$ mv bash_history.new ~/.bash_history
$ history -c  # Clear current session's history
$ history -n  # Read history from history file
```

### Development Setup

1. Clone the repo: `git clone https://github.com/n8henrie/fixbashhistory && cd
   fixbashhistory`

## Troubleshooting / FAQ

- What are some cool history settings for bash >= 4?
    - `man bash` and search for `shopt`
    - Google: `histappend cmdhist lithist histreedit histverify`
