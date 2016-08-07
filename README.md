# alfred
Because even Batman needs a little help.

![Alfred](https://raw.githubusercontent.com/kcmerrill/alfred/master/alfred.jpg "Alfred")

## Installation
``` $ go get github.com/kcmerrill/alfred ```

## What is it
A simple go/yaml powered make file with a bit of a twist.

## Features
* Tasks can call other tasks
* Tasks can be run every _insert a proper go duration_
* Tasks can pause
* Tasks can call other tasks depending on success/failure
* Alfred uses go templates so you can inject variables making tasks reusable
* Common tasks can be invoked inside using alfred
* Common tasks will be stored in this repository for shared use(git, docker, mail, slack, etc ...) (coming soon)
* Optional private/public repositories so you can share private tasks with coworkers (coming soon)
* Start alfred as a webserver to start tasks remotely
* No need to be in the same directory when running alfred if it's local, as long as it's in a parent directory
* Fire off multiple tasks at once (coming soon)


## Why
I have a lot of tasks I do daily that I'd like a remote repository to use. Not only can I update them, but others can also use it too. I love make files, but wanted just a bit more functionality while keeping it simple. I also love ansible but it's a bit too heavy handed for what I need, so this is my attempt at KISS.

## Example uses
* Monitor webistes with reusable tasks(see example)
* Setup/Update/Deploy projects in your dev env
* Start/Stop remote tasks
* Simple Nagios, Jenkins, pingdom replacement

## Quick docs
Using the example below there a few things to notice.

1. `check` `mail` `slack` `down` `up` `notify` are all names of tasks. `check` is the main entrypoint task, however, you can put whatever you'd like for your own use case. You can call each one via `alfred _taskname_`.
2. `usage` gives a quick example of how to use that particular task
3. `command` is a string, or in this case the `|` denotes a multiline string to be run. It's run via `bash -c`
4. If `command` exits properly with a `0` exit code, call tasks within `ok`
5. `ok`, `fail`, `tasks` all can take a string with space seperated task names. Note the `notify` task.
5. `wait` a simple pause. You can use sleep X however, wait takes in a string that parses into a go duration. So `1s` = 1 second, `1h` = 1 hour,  `1m` = 1 month and so on.
6. `summary` explains what is happening when it's happening
7. `every` is also a go duration(see `wait`). That means run this task(and all of it's tasks that it calls) every so often. Useful for checking things.
8. `wait` and `command` uses go text/templates. The main thing that provides at this point is access to `.Args`. Note, index .Args 0 starts immediately after the task to run.

## Example
```  check:
      summary: Checks a page
      usage: alfred check personal http://kcmerrill.com kcmerrill@gmail.com 10m
      command: |
          mkdir -p /tmp/website/checks
          wget {{ index .Args 1 }} -O /tmp/website/checks/{{ index .Args 0}}
      ok: up
      fail: notify
      wait: 3s
      every: "{{ index .Args 3 }}"

  mail:
      summary: Sending an email
      command: |
          echo "{{ index .Args 0 }} site is down! " |  mail -s "Website down" {{ index .Args 2}}

  slack:
      summary: "Sending a slack message to #general"

  down:
      summary: The website is down

  up:
      summary: The website is up

  notify:
      summary: Notify the website is down
      tasks: down slack mail```
