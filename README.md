# alfred
Because even Batman needs a little help.

[![Build Status](https://travis-ci.org/kcmerrill/alfred.svg?branch=master)](https://travis-ci.org/kcmerrill/alfred) [![Join the chat at https://gitter.im/kcmerrill/alfred](https://badges.gitter.im/kcmerrill/alfred.svg)](https://gitter.im/kcmerrill/alfred?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

![Alfred](https://raw.githubusercontent.com/kcmerrill/alfred/master/assets/alfred.jpg "Alfred")

## What is it
A simple go/yaml powered make file/task runner with a bit of a twist.

## Features
- Extendable. Common tasks(Private too)
- Watch files for modifications
- Retry/Rerun tasks based on failures before giving up 
- Logging
- Success/Failure decision tree
- Run tasks asynchronously or synchronously 
- Autocomplete task names
- Static webserver
- Many more! 

## Usage
Create a file named: `alfred.yml`
```
say.hello:
    summary: I will say hello!
    usage: alfred say.hello
    command: echo "Hello!"

say.howareyou:
    summary: I will ask how you are
    usage: alfred say.howareyou
    command: |
        echo "How"
        echo "Are"
        echo "You?"

speak:
    tasks: say.hello say.howareyou

blurt:
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

[For additional documentation, please refer to the GUIDE](GUIDE.md "additional documentation")

## Example uses
* Monitor websites
* Setup/Update/Deploy projects in your dev env
* Simple Nagios, Jenkins, pingdom replacement
* Monitor crons(alert on failures, update endpoints etc ... )
* Watch for file modifications to run tests->builds


## Alfred files getting too large?
You can break up your alfred files in multiple ways. The following are glob patterns that can be used:`/alfred.yml`, `/.alfred/*alfred.yml`, `/alfred/*alfred.yml`. As an example, you can create a directory called `alfred` or `.alfred` or just create mutliple alfred files.

## Tab completion

Copy the included `alfred.completion.sh` to `/etc/bash_completion.d/`, or source it in your `~/.profile` file.

## Testing
You might say I've cheated the testing route by only scraping the output. You'd be right.

"I live with a wizard. I cheat" ~ Mouse
