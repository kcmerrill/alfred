# alfred
Because even Batman needs a little help.

[![Build Status](https://travis-ci.org/kcmerrill/alfred.svg?branch=master)](https://travis-ci.org/kcmerrill/alfred) [![Join the chat at https://gitter.im/kcmerrill/alfred](https://badges.gitter.im/kcmerrill/alfred.svg)](https://gitter.im/kcmerrill/alfred?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

![Alfred](https://raw.githubusercontent.com/kcmerrill/alfred/master/assets/alfred.jpg "Alfred")

## Installation
``` $ go get github.com/kcmerrill/alfred ```

Fan of docker? It's in a container too!

```$ alias alfred="docker run -ti --rm -v $PWD:$PWD -w $PWD kcmerrill/alfred"```


## What is it
A simple go/yaml powered make file with a bit of a twist.

## Features
* Tasks can call other tasks
* Tasks can be run every `time.Duration`
* Tasks can pause or wait for an event to happen
* Tasks can call other tasks depending on success/failure
* Alfred uses go templates so you can inject variables making tasks reusable
* Common tasks can be invoked inside using alfred
* Common tasks will be stored in this repository for shared use(git, docker, mail, slack, etc ...) (Available ... more alwyas being added)
* Optional private/public repositories so you can share private tasks with coworkers
* No need to be in the same directory when running alfred if it's local, as long as it's in a parent directory
* Fire off multiple tasks at once
* Variables, regular text and command line evaluated


## Why
I have a lot of tasks I do daily that I'd like a remote repository to use. Not only can I update them, but others can also use it too. I love make files, but wanted just a bit more functionality while keeping it simple. I also love ansible but it's a bit too heavy handed for what I need, so this is my attempt at KISS.

## Example uses
* Monitor webistes with reusable tasks(see example)
* Setup/Update/Deploy projects in your dev env
* Start/Stop remote tasks
* Simple Nagios, Jenkins, pingdom replacement
* Monitor crons(alert on failures, update endpoints etc ... )

## Screencast
A video(~35 minutes long) showing alfred and how to use it. Using contrivied examples, I believe it should get the point across.

[![Alfred Overview](http://i3.ytimg.com/vi/v2ivtM5anbk/hqdefault.jpg)](https://www.youtube.com/watch?v=v2ivtM5anbk)

## Docker-compose replacement(as an example)
I recently used alfred to setup a process that cuts the time to download/lift/build containers literally in half.

![Alfred](https://raw.githubusercontent.com/kcmerrill/alfred/master/assets/alfred_benchmark.png "Alfred")

## Example 1
Use alfred to monitor crons. Just wrap your command in an alfred task, and depending on success failure do something with it. In this case, we can use alfred to send in a datapoint to datadog

```
* * 1 * * alfred monitor python somescript.py
```

```
monitor:
    summary: Monitor a specific cron job
    command: {{ .AllArgs }}
    ok: success
    fail: failed

ok:
    summary: Send ok data point to datadog agent
    command: |
        echo "cron.{{ index .Args 1 }}.ok:1|c|#cron" | nc -w 1 -u 0.0.0.0 8125
    private: true

failed:
    summary: Send failure data point to datadog agent
    command: |
        echo "cron.{{ index .Args 1 }}.failed:1|c|#cron" | nc -w 1 -u 0.0.0.0 8125
    private: true

```

## Example 2
This example demonstrates the reuseability of alfred. This check is hosted within the common module `check`, so as long as you have the binary installed, you can tap into quite a few shared libraries. This is one of them. This particular module is dependant on another module, the `notify` module.

In order to use this, common module simply run the following command, replacing `kcmerrill.com` with your website, and your `supersecretkey` with a slack incoming webhook key/secret(typically the last two segments of the webhook url)

`alfred /http.slack "kcmerrill.com" "supersecretkey"` *10

`* optional int in minutes to notify intervals when the website is down`

By leveraging the common module, alfred will post into slack letting you know that your website is down.

```
http.slack:
    summary: HTTP Check
    usage: alfred /http.slack "website" "supersecretkey"
    dir: /tmp/http/{{ index .Args 0 }}
    command: wget {{ index .Args 0 }}
    ok: cleanup
    fail: send.notification
    every: 10s

send.notification:
    dir: /tmp/http/{{ index .Args 0 }}
    command: test ! -f last_notified || test $(find -amin +{{ index .Args 2 }} | wc -l) -ge "2"
    private: true
    ok: notify
    defaults:
        - ""
        - ""
        - "10"

notify:
    dir: /tmp/http/{{ index .Args 0 }}
    summary: Send slack notification
    notify: slack "{{ index .Args 1 }}" "{{ index .Args 0 }} is down"
    command: touch last_notified

cleanup:
    dir: /tmp/http/{{ index .Args 0 }}
    command: |
        rm -rf last_notified
        rm -rf index.html*
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
    fail: ten

eight:
    summary: This was only called because step seven was succesful.
    command: echo "Also,notice this is private. Cannot be called directly"
    private: true

nine:
    summary: Try to ls a folder that _hopefully_ doesn't exist. Notice exit code
    command: echo "Even though alfred worked, you can specifically set exit codes." && ls /kcwashere
    fail: ten
    ok: eight

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
    summary: Run ls every 1 seconds, or any golang duration
    command: ls
    every: 1s

fourteen:
    summary: Pass along arguments using go/text templates. Try running without an argument.
    command: ls {{ index .Args 0 }}
    usage: alfred fourteen foldername

fifteen:
    summary: Pass along arguments again ... but use defaults
    command: ls {{ index .Args 0 }}
    defaults:
        - "."

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
        mkdir directory
        cd directory
        echo "Current directory:"
        pwd
        echo "Now call alfred four, which is an alfred command that pwd"
        alfred four

twentyone:
    summary: Even though alfred worked, If the command failed you can still exit after running all of the failed tasks
    command: ls /asdf
    fail: four
    exit: 42

twentytwo:
    summary: Multitasks! You can combine this with other things too!
    multitask: long-task1 long-task2 long-task1 long-task2

alfred.vars:
    firstname: casey
    user: whoami
    pwd: pwd

twentythree:
    summary: Lets test out alfred.vars
    command: |
        echo The variable firstname lastname = {{ index .Vars "firstname" }} {{ index .Vars "lastname" }}
        echo You are the user {{ index .Vars "user" }}
        echo The pwd of this directory is {{ index .Vars "pwd" }}
    #Default vars if none are set ...
    vars:
        firstname: kc
        lastname: merrill

long-task1:
    summary: This long task takes 10 seconds
    command: |
        sleep 10
        echo "Done with 10 second long task"

long-task2:
    summary: This long task takes 9 seconds
    command: |
        sleep 9
        echo "Done with 9 second long task"

twentyfour:
    summary: Lets test out the log!
    log: /tmp/log.txt
    command: echo "This should be in /tmp/log.txt"

twentyfive:
    summary: Skip! Same as "exit" however, it only skips the task and _not_ exit the application
    command: ls /shouldhopefullyfail
    skip: true


twentysix*:
    summary: Notice the astrick? It means it's a "main" task. Useful for a long alfred file
    command: echo "See how it pops out from the rest?"
```

## Remote/Custom modules
In order to extend alfred to things out of the box, you can create your own modules. To do so, create an alfred configuration file. Here is an example:
```
$ cat ~/.alfred/config.yml
remote:
    kcmerrill: http://kcmerrill.com:8081/shares/
```
You can add as many remotes as you'd like. By default there will be one remote automagically added. `common`. A remote consist of a name, a forward slash and a name. If you're ok with having your custom work shared on github, you can setup a module repository.

Alfred comes with a really basic web server so you can host private/sensative modules on your internal network. To start the webserver you can simply: `alfred --serve --dir . --port 8080`. Note, `dir` and `port` are not required and default to `.` and `8080` respectively.

The folder in which you start serving your alfred files should contain a `modulename/alfred.yml` and inside the alfred.yml is your standard yaml file.

## Common modules
![Alfred](https://raw.githubusercontent.com/kcmerrill/alfred/master/assets/alfred_slack.png "Alfred")

This is the `notify` module in action.
