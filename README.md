# alfred
Because even Batman needs a little help.

[![Build Status](https://travis-ci.org/kcmerrill/alfred.svg?branch=master)](https://travis-ci.org/kcmerrill/alfred) [![Join the chat at https://gitter.im/kcmerrill/alfred](https://badges.gitter.im/kcmerrill/alfred.svg)](https://gitter.im/kcmerrill/alfred?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

![Alfred](https://raw.githubusercontent.com/kcmerrill/alfred/master/assets/alfred.jpg "Alfred")

## What is it
A simple go/yaml powered make file/task runner with a bit of a twist.

## Usage
Create a file named: `alfred.yml`
```

# Create a task, name it whatever you'd like.
say.hello:
    # Lets give it a quick summary. Optional.
    summary: I will say hello!
    # Describe how to use this task. Optional.
    usage: alfred say.hello
    # The command to perform
    command: echo "Hello!"

say.howareyou:
    # Lets give it a quick summary. Optional.
    summary: I will ask how you are
    # Describe how to use this task. Optional.
    usage: alfred say.howareyou
    # The command to perform
    command: |
        echo "How"
        echo "Are"
        echo "You?"

speak:
    # You can call multiple tasks in an order
    tasks: say.hello say.howareyou

blurt:
    # You can run multiple tasks at the same time
    multitask: say.hello say.howareyou

```

Then, anywhere in the top-level or child directories to the `alfred.yml` file:

`alfred` Will show you all of the available tasks and a quick summary.

`alfred say.hello` Will simply say hello

`alfred say.howareyou` Will ask how you are

`alfred speak` will perform both tasks in the specified order

`alfred blurt` will perform both tasks at the same time

# Quick Walkthrough
[![asciicast](https://asciinema.org/a/103711.png)](https://asciinema.org/a/103711)

# Every option
### Note, you do not need to use them all, however, you can.

Lets create a task that has _every_ option available in the -order- it's run in(note, you can put them in any order, but they will be executed in the following order)

```
alfred.vars:
    var.one: somevar
every.option*: # can be named anything you want. An `*` denotes it's an "important" task
    setup: task.one task.two task.three # space separated list of task names. Run first before anything else
    alias: every option # string separated list fof aliases for this task.
    dir: /tmp # defaults to where alfred.yml is found, else, uses this option. Dir is created if not exist
    log: /tmp/log_output.txt # a log where all stdout/stdin of `command` is stored
    watch: '.*\.go' # a regular expression, that will watch for any files changed within the last second matching regex
    modulenamehere: docker kill.remove containername # Anything that is not a valid key is a module(a task that is defined remotely)
    summary: A quick description of this task.
    retry: 3 # How many times we should attempt to run the command option before giving up
```

## Example uses
* Monitor webistes with reusable tasks(see example)
* Setup/Update/Deploy projects in your dev env
* Start/Stop remote tasks
* Simple Nagios, Jenkins, pingdom replacement
* Monitor crons(alert on failures, update endpoints etc ... )
* Watch for file modifications to run tests->builds

## Common modules
![Alfred](https://raw.githubusercontent.com/kcmerrill/alfred/master/assets/alfred_slack.png "Alfred")

This is the `notify` module in action. `alfred /notify slack ...`

## Alfred files getting too large?
You can break up your alfred files in multiple ways. The following are glob patterns that can be used:`/alfred.yml`, `/.alfred/*alfred.yml`, `/alfred/*alfred.yml`. As an example, you can create a directory called `alfred` or `.alfred` or just create mutliple alfred files.

## Tab completion

Copy the included `alfred.completion.sh` to `/etc/bash_completion.d/`, or source it in your `~/.profile` file.

## Testing
You might say I've cheated the testing route by only scraping the output. You'd be right.

"I live with a wizard. I cheat" ~ Mouse
