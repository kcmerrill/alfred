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
* Common tasks will be stored in this repository for shared use(git, docker, mail, slack, etc ...) (Available ... more coming soon)
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

## Example 1
```
check:
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
      tasks: down slack mail
```

## Example 2
You can use alfred to get a project running. Useful if your projects have a bajillion steps, or if you're like me and you are typically responsible for dev enviornments at your work. Using alfred, you can run setup. In this example, we will use a common github module to clone a github project in a folder of your choosing, update all the submodules, create symlinks composer update and run version as the final check to ensure things are working.

```alfred common/github setup kcmerrill/yoda yoda```

We used common/github to clone the repository _FIRST_ then run alfred setup inside of it. This is useful because some projects are private and you do not have access to the alfred file like you normally would.

If it is a public repository, like `kcmerrill/yoda` is, then you can simply call it remotely. In this case, `install` is the task that will get yoda setup.

```alfred kcmerrill/yoda install```

That will seek out the project on github, look in the master branch and then call the `install` task within the yaml file. Which will then proceed to check out the code for you. Take a peek at `kcmerrill/yoda` for it's `alfred.yml` to see how it's setup.

## Example 3 (demo-everything)
You can see this file in the examples folder. I will try to update this when features get added

```
one:
    summary: Displaying the task name

two:
    summary: A simple echo
    command: echo "A simple echo command"

three:
    summary: Change the working directory
    dir: /tmp
    command: pwd

four:
    summary: Notice how the directory has changed back?
    command: pwd

five:
    summary: Step five, but aliased as step six too! Space seperated
    command: ls
    alias: six six.one six.two

seven:
    summary: Step seven is a simple ls, but will automagically call step 8
    command: ls
    ok: eight

eight:
    summary: This was only called because step seven was succesful.
    command: echo "Also,notice this is private. Cannot be called directly"
    private: true

nine:
    summary: Try to ls a folder that _hopefully_ doesn't exist. Notice exit code
    command: echo "Even though alfred worked, you can specifically set exit codes." && ls /kcwashere
    fail: ten

ten:
    summary: Only called because step 9 failed
    command: echo "Again, notice how this step is private"
    private: true

eleven:
    summary: Call multiple tasks as a task group, space seperated
    tasks: four five six

twelve:
    summary: Call alfred within itself
    command: alfred eleven

thirteen:
    summary: Run ls every 3 seconds, or any golang duration
    command: ls
    every: 3s

fourteen:
    summary: Pass along arguments using go/text templates. Try running without an argument.
    command: ls {{ index .Args 0 }}
    usage: alfred fourteen foldername

fifteen:
    summary: Pass along arguments again ... but use defaults
    command: ls {{ index .Args 0 }}
    defaults:
        - /tmp

sixteen:
    summary: Remotes allow you to reuse common components. This will completely setup a git project as an example
    dir: /tmp
    git: clone kcmerrill/alfred alfred

seventeen:
    summary: Wait! Sure, you can sleep, but this will let you do so via a golang duration
    wait: 5s
    command: ls

eighteen:
    summary: You can combine everything you've seen above. Infinite loop
    command: test $(whoami) = "rooto"
    wait: 10s
    every: 7s
    fail: nineteen
    ok: nineteen.1
    exit: 42

nineteen:
    summary: You are not root, but checkout this multiline command while you're here
    private: true
    command: |
        cd /tmp && pwd
        cd /tmp
        pwd

nineteen.1:
    summary: You apparently _ARE_ root. Cowers in feer of your l33t/\/355.
    private: true
    command: |
        echo "Checkout this multi line command!"
        cd /tmp && pwd
        cd /tmp
        pwd

twenty:
    summary: As long as an alfred file is in a parent directory, you can call it and alfred will find it
    command: |
        mkdir fakedirectory
        cd fakedirectory
        echo "Current directory:"
        pwd
        echo "Now call alfred four, which is an alfred command that pwd"
        alfred four

twentyone:
    summary: Even though alfred worked, If the command failed you can still exit after running all of the failed tasks
    command: ls /asdf
    fail: four
    exit: 42
```
