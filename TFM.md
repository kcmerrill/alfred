# The Manual

Alfred tasks can be complex or as simple as you want them to be. The idea is they are small, reusable and can help automate your daily workflow. Each task is made up of [Components](#components). You can think of these components as building blocks to build whatever it is you need. In this manual, we'll go over some classic use cases for alfred, we'll discuss each component in depth, tips and tricks and also ways of structuring your tasks for maximum awesomeness.

If given enough building blocks anything is possible, so alfred really is up to you to choose your own adventure. There are some plugins, and remote tasks available to you, but the power behind alfred is it's flexibility to fit into your specific usecases. I'd like to think so anyways.

* [Usage](#usage)
  * [Tasks](#tasks)
    * [Local tasks](#local-tasks)
    * [Remote tasks](#remote-tasks)
  * [Arguments](#arguments)
  * [Taskgroups](#taskgroups)
  * [Golang Templating](#golang-templating)
    * [Vars](#vars)
    * [Stdin](#stdin)
  * [Components](#components)
    * [log | string](#log--string)
    * [defaults | []string](#defaults--string)
    * [summary | string](#summary--string)
    * [string | string(text, command)](#stdin--stringtext-command) 
    * [dir | string(dir, command)](#dir--stringdir-command)
    * [config | string(filename, yaml)](#config--stringfilename-yaml)
    * [prompt | map[string]string](#prompt--mapstringstring)
    * [register | map[string]string](#register--mapstringstring)
    * [env | map[string]string](#env--mapstringstring)
    * [serve | string(port)](#serve--stringport)
    * [setup | TaskGroup\{\}](#setup--taskgroup)
    * [multitask | TaskGroup\{\}](#multitask--taskgroup)
    * [tasks | TaskGroup\{\}](#tasks--taskgroup)
    * [watch | string(regExp)](#watch--stringregexp)
    * [for | map[string]string](#for--mapstringstring)
      * [tasks | TaskGroup\{\}](#tasks--taskgroup-1)
      * [multitask | TaskGroup\{\}](#multitask--taskgroup-1)
      * [args | string(command, text)](#args--stringcommand-text)
    * [command | string](#command--string)
    * [commands | string](#commands--string)
    * [ok | TaskGroup\{\}](#ok--taskgroup)
    * [http\.tasks | map[string]string](#httptasks--mapstringstring)
    * [fail | TaskGroup\{\}](#fail--taskgroup)
    * [wait | duration](#wait--duration)
    * [every | duration](#every--duration)
* [Tips and Tricks](#tips-and-tricks)
    * [Aliases](#aliases)
    * [Task Inheritance](#task-inheritance)
    * [Run a task X times](#run-a-task-x-times)

# Usage

There are a few underlying key concepts throughout Alfred that should be pointed out. [Taskgroups](#taskgroups), [Components](#components), [Arguments](#arguments). We'll go over a few of them here.

To see a list of tasks, simply invoke `alfred`. Any task that is deemed important by the tasks creator(read: any task that has a summary) will be displayed here.

## Tasks

### Local tasks

By default, alfred will look in the current directory for an `alfred.yml` file. If it cannot be found, it will continue to go up the directory structure until it reaches the root directry. At this point, if the task provided does not exist, it will exit. 

If your alfred files are getting large, you can break up your files by creating a `.alfred` or `alfred` folder, and inside create new .yml files named `something.alfred.yml`. Each file is then concatonated together, so be sure you do not have any task name collisions. 

All tasks start where the alfred file/folder are located(unless it's a remote task)

### Remote tasks

One bit of functionality that makes alfred so flexible is the ability to have private/remote repositories of alfred files. While this feature is not new to alfred, the way it's invoked is. By using a `/` at the start of the task name, alfred knows to lookup tasks in the newly created repository [kcmerrill/alfred-tasks](https://github.com/kcmerrill/alfred-tasks)

Also, you can use a web address in order to get access to tasks. Of course, if the web address is private, then your tasks are protected.

Notice the `:` which distinguishes the URL from the taskname. Without a taskname, alfred will list the available tasks.

```sh
07:52 PM ✔ kcmerrill (v0.2) alfred ] alfred /testing:tdd.go
[ 0s] (26 Dec 17 19:52 MST) tdd.go started [] Watch .go files and run test.go
[ 0s] (26 Dec 17 19:52 MST) tdd.go watching ./
```

```sh
07:54 PM ✘ kcmerrill (v0.2) alfred ] alfred https://raw.githubusercontent.com/kcmerrill/alfred-tasks/master/testing.yml:tdd.go
[ 0s] (26 Dec 17 19:54 MST) tdd.go started [] Watch .go files and run test.go
[ 0s] (26 Dec 17 19:54 MST) tdd.go watching ./
```

## Arguments

In order to make tasks reusable, you can pass in arguments. This is true for the yaml file configuration along with the CLI. The arguments are positional, so the very first argument after the task to run is location `0`. The next is `1` and so on. You can obtain these variables by using golang templating by putting in the representation. 

```yaml
task.name:
    summary: Display the arguments
    command: |
        echo The first argument is: {{ index .Args 0 }}
        echo All of the arguments passed in: {{ .AllArgs }}
```
```sh
07:19 PM ✔ kcmerrill (v0.2) demo ] alfred task.name one two three
[ 0s] (26 Dec 17 19:19 MST) task.name started [one, two, three] Display the arguments
[ 0s] (26 Dec 17 19:19 MST) The first argument is: one
[ 0s] (26 Dec 17 19:19 MST) All of the arguments passed in: one two three
[ 0s] (26 Dec 17 19:19 MST) task.name ✔ ok [one, two, three] elapsed time '0s'
07:19 PM ✔ kcmerrill (v0.2) demo ]
```

You will get an error if the template is unable to find the appropriate arguments

## Taskgroups

There are several components that call other tasks. These are called TaskGroups. A few example of these components would be [tasks](#tasks--taskgroup), [multitask](#multitask--taskgroup), [setup](#setup--taskgroup) to name a few. 

You can define a task group in multiple ways. For task groups that do not need to call paramaters, you can simply put them on a single space delimited line. 

If your task requires arguments, or a mix of no arguments and arguments, put the tasks on new lines. This way you can mix and match the type of tasks that need and don't need arguments. This also helps task readability.

```yaml
task.one:
    command: echo task.one

task.two:
    command: echo task.two

task.three:
    command: echo {{ index .Args 0 }}
    defaults:
        - "arg0"
        - "arg1"
        - "arg2"

taskgroup.inline:
    setup: task.one task.two task.three
    command: echo all finished

taskgroup.mixed:
    setup: |
        task.one
        task.three({{ index .Args 0 }}, second.arg, third.arg)
        task.two
```


###

## Golang Templating

Another flexibile feature of Alfred is the ability to use the go templating language within your yaml files. As demonstrated through the documentation, you can do extra complex things based on this. You can read more about [golang templates here](https://golang.org/pkg/text/template/). Included with the templates is [masterminds/sprig](http://masterminds.github.io/sprig/) which include a ton of extra handy functionality. Setting defaults, uuids, env functions, Date, function among a whole host of other awesome goodies. 

### Vars

You can register vars many different ways through various components. You can access these vars via the `.Vars` object in your templates. `{{ .Vars.varname }}`

### Stdin

Stdin can be accessed through the variable `{{ .Stdin }}`. This can be handy if you need to pipe text into commands. You can also loop over standard in by giving it to the `for` component for example. Take a peek at the [stdin](#stdin) component for more functionality.

A quick usecase.

```yaml
superheros:
    summary: Show the superheros!
    for:
        args: "{{ .Stdin }}"
        tasks: echo.superhero

echo.superhero:
    summary: Echo the superhero name
    command: |
        echo "Superhero {{ index .Args 0 }}"
```

```sh
02:48 PM ✔ kcmerrill  tmp ] cat superheros.txt
superman
batman
spiderman
02:48 PM ✔ kcmerrill  tmp ] cat superheros.txt | alfred superheros
[    0s] (27 Dec 17 14:48 MST) superheros started [] Show the superheros!
[    0s] (27 Dec 17 14:48 MST) echo.superhero started [superman] Echo the superhero name
[    0s] (27 Dec 17 14:48 MST) Superhero superman
[    0s] (27 Dec 17 14:48 MST) echo.superhero ✔ ok [superman] elapsed time '0s'
[    0s] (27 Dec 17 14:48 MST) echo.superhero started [batman] Echo the superhero name
[    0s] (27 Dec 17 14:48 MST) Superhero batman
[    0s] (27 Dec 17 14:48 MST) echo.superhero ✔ ok [batman] elapsed time '0s'
[    0s] (27 Dec 17 14:48 MST) echo.superhero started [spiderman] Echo the superhero name
[    0s] (27 Dec 17 14:48 MST) Superhero spiderman
[    0s] (27 Dec 17 14:48 MST) echo.superhero ✔ ok [spiderman] elapsed time '0s'
[    0s] (27 Dec 17 14:48 MST) superheros ✔ ok [] elapsed time '0s'
02:48 PM ✔ kcmerrill  tmp ]
```

## Components 

The components here will be listed in order in which they are executed within Alfred. With the way golang's maps work, [they are randomized](https://github.com/golang/go/issues/2630) to prevent DOS attacks. The reason this is important? Your components within your tasks can be ordered however you'd like, but they will be executed in a specific order. 


### log | string

The raw output of commands will be sent to the specified filename appended. If the file does not exist it will attempt to be created.

A quick note, multiple tasks can call different logs, which means that any parent tasks will also have the output sent to said log file. 

```yaml
log.task:
  summary: Demo of the 'log' component
  log: /tmp/log.demo.txt
  command: |
    echo "This should be in the log /tmp/log.demo.txt"
```

```sh
09:09 PM ✔ kcmerrill (v0.2) demo ] alfred log
[ 0s] (25 Dec 17 21:09 MST) log started [] Demo of the 'log' component
[ 0s] (25 Dec 17 21:09 MST) This should be in the log /tmp/log.demo.txt
[ 0s] (25 Dec 17 21:09 MST) log ✔ ok [] elapsed time '0s'
09:09 PM ✔ kcmerrill (v0.2) demo ] cat /tmp/log.demo.txt
This should be in the log /tmp/log.demo.txt
09:09 PM ✔ kcmerrill (v0.2) demo ]
```

### defaults | []string

When passing in arguments, it's possible to set default arguments. If a given value is an empty string, it will still be required. A use case would be accepting three arguments. This is so you can set the default value to the third argument, making the first two required. 

```yaml
arguments.task:
    summary: Demo of arguments
    command: |
        echo {{ index .Args 0 }}
        echo {{ index .Args 1 }}
        echo {{ index .Args 2 }}
    defaults:
        - "" #arg0 is required
        - "" #arg1 is required
        - I am a default
```

```sh
09:17 PM ✔ kcmerrill (v0.2) demo ] alfred arguments.task one
[ 0s] (25 Dec 17 21:17 MST) arguments.task template Invalid Argument(s)
[ 0s] (25 Dec 17 21:17 MST) arguments.task ✘ failed [one] elapsed time '0s'
09:17 PM ✘ kcmerrill (v0.2) demo ] alfred arguments.task one two
[ 0s] (25 Dec 17 21:17 MST) arguments.task started [one, two, I am a default] Demo of arguments
[ 0s] (25 Dec 17 21:17 MST) one
[ 0s] (25 Dec 17 21:17 MST) two
[ 0s] (25 Dec 17 21:17 MST) I am a default
[ 0s] (25 Dec 17 21:17 MST) arguments.task ✔ ok [one, two, I am a default] elapsed time '0s'
09:17 PM ✔ kcmerrill (v0.2) demo ]
```

### summary | string

The task can be described by it's summary. When performing a list of the tasks, or when the task is started, the summary will be displayed.

```yaml
show.summary:
  summary: This is the summary. 
```

```sh
09:17 PM ✔ kcmerrill (v0.2) demo ] alfred show.summary
[ 0s] (25 Dec 17 21:21 MST) show.summary started [] This is the summary.
[ 0s] (25 Dec 17 21:21 MST) show.summary ✔ ok [] elapsed time '0s'
09:21 PM ✔ kcmerrill (v0.2) demo ]
```

### stdin | string(text, command)

Using the stdin component will allow you to chain tasks. `stdin` component will be evaulated as a command, and if a non zero code is returned, the result of the string will be stored in stdin, and will then be sent to each subsequent command down stream. If the text of `stdin` is not a valid command, resulting in a non zero exit code, will be left as regular text and sent to eachh subsequent command down stream. 

```yaml
pipe:
    summary: Add to pipe
    stdin: "helloworld"
    ok: md5

md5:
    command: |
        md5
```

```sh
11:07 PM ✔ kcmerrill  tmp ] alfred pipe
[    0s] (28 Dec 17 23:07 MST) pipe started [] Add to pipe
[    0s] (28 Dec 17 23:07 MST) pipe ✔ ok [] elapsed time '0s'
[    0s] (28 Dec 17 23:07 MST) pipe ok.tasks md5
[    0s] (28 Dec 17 23:07 MST) md5 started []
[    0s] (28 Dec 17 23:07 MST) fc5e038d38a57032085441e7fe7010b0
[    0s] (28 Dec 17 23:07 MST) md5 ✔ ok [] elapsed time '0s'
11:07 PM ✔ kcmerrill  tmp ] printf "helloworld" | alfred md5
[    0s] (28 Dec 17 23:07 MST) md5 started []
[    0s] (28 Dec 17 23:07 MST) fc5e038d38a57032085441e7fe7010b0
[    0s] (28 Dec 17 23:07 MST) md5 ✔ ok [] elapsed time '0s'
```

### dir | string(dir, command)

Set the directory to where the task should be run. Any component or command run will be relative to `dir`. The string will be `evaluated` and if a valid exit code is returned, will be the value set. By default, the directory is set to where the alfred files are found.

One thing to note. Be careful when using dir with `multitask`, as the order is not guarenteed to be run.

```yaml
dir:
    summary: Lets display the directory
    dir: "{{ index .Args 0 }}"
    command: pwd
```

```sh
09:59 PM ✔ kcmerrill (v0.2) demo ] alfred dir
[ 0s] (25 Dec 17 21:59 MST) dir started [] Lets display the directory
[ 0s] (25 Dec 17 21:59 MST) /private/tmp/does/not/exist
[ 0s] (25 Dec 17 21:59 MST) dir ✔ ok [] elapsed time '0s'
09:59 PM ✔ kcmerrill (v0.2) demo ]
```

### config | string(filename, yaml)

A valid filepath that contains a `yaml` unmarshable string of `key: value` pairs. These pairs will then be set as a [variable](#variable). If the filepath does not exist, then the string set itself will be unmarshed as `key: value` pairs and will also be set as variables. The values will then be `evaluated` as a command, and if a valid exit code of zero is returned, the `CombinedOutput()` will be the new value. 

```yaml
configuration:
    summary: This task will show how to use config
    config: |
        user: whoami
        email: "kcmerrill@gmail.com"
    command: |
        echo "The current user is {{ .Vars.user }}"
        echo "The current user's email address is {{ .Vars.email }}"
        echo "The user twitter handle is {{ default "is not set" .Vars.twitter }}"
    
configuration.file:
    summary: This will show a configuration with a valid yaml file
    config: /tmp/file.that.exists.yml
    command: |
        echo "The current user is {{ .Vars.user }}"
        echo "The current user's email address is {{ .Vars.email }}"
        echo "The user twitter handle is {{ default "is not set" .Vars.twitter }}"
    

```

```sh
09:34 PM ✔ kcmerrill (v0.2) demo ] alfred configuration
[ 0s] (25 Dec 17 21:34 MST) configuration started [] This task will show how to use config
[ 0s] (25 Dec 17 21:34 MST) The current user is kcmerrill
[ 0s] (25 Dec 17 21:34 MST) The current user's email address is kcmerrill@gmail.com
[ 0s] (25 Dec 17 21:34 MST) The user twitter handle is is not set
[ 0s] (25 Dec 17 21:34 MST) configuration ✔ ok [] elapsed time '0s'
```

### prompt | map[string]string

A `key:value` pair, where the key is the registered variable, and the value is the phrase that the user will be prompted for.

The values can then be access via `{{ .Vars.<key>}}`

```yaml
prompt:
    summary: lets ask some questions
    prompt:
        fav: What is your favorite color? 
        car: Whatis your favorite car?
    command: |
        echo Your favorite car is {{ .Vars.car }}
        echo Your favorite color is {{ .Vars.fav }}
```

```sh
02:05 PM ✔ kcmerrill (master) alfred ] alfred prompt
[    0s] (27 Dec 17 14:05 MST) prompt started [] lets ask some questions
[    0s] (27 Dec 17 14:05 MST) prompt prompt What is your favorite color? asdf
[    1s] (27 Dec 17 14:05 MST) prompt prompt Whatis your favorite car? fff
[    3s] (27 Dec 17 14:05 MST) prompt registered car fff
[    3s] (27 Dec 17 14:05 MST) prompt registered fav asdf
[    3s] (27 Dec 17 14:05 MST) Your favorite car is fff
[    3s] (27 Dec 17 14:05 MST) Your favorite color is asdf
[    3s] (27 Dec 17 14:05 MST) prompt ✔ ok [] elapsed time '3s'
02:05 PM ✔ kcmerrill (master) alfred ]
```

### register | map[string]string 

Based on `key: value` pairs, will register the pairs as [variables](#variables). The value is then `evaluated` when a zero exit code is shown, the `CombinedOutput()` is the resulting value. 

```yaml
register:
    summary: Demonstrate the registration of variables
    register:
        user: whoami
        twitter: "@themayor"
    command: |
      echo "{{ index .Vars "user" }}"
      echo "{{ .Vars.twitter }}"
```

```sh
09:43 PM ✔ kcmerrill (v0.2) demo ] alfred register
[ 0s] (25 Dec 17 21:43 MST) register started [] Demonstrate the registration of variables
[ 0s] (25 Dec 17 21:43 MST) register registered user kcmerrill
[ 0s] (25 Dec 17 21:43 MST) register registered twitter @themayor
[ 0s] (25 Dec 17 21:43 MST) kcmerrill
[ 0s] (25 Dec 17 21:43 MST) @themayor
[ 0s] (25 Dec 17 21:43 MST) register ✔ ok [] elapsed time '0s'
09:43 PM ✔ kcmerrill (v0.2) demo ]
```

### env | map[string]string

Setting env variables is a lot like the `register` component, and the `config` component as well, however the difference is the variables will be available as ENV variables on the CLI. 

```yaml
env:
  summary: Lets set some env variables!
  env: 
    WHO: whoami
    TWITTER: "{{ index .Args 0 }}"
  command: |
    echo twitter:$TWITTER
    echo who:$WHO
```

```sh
09:53 PM ✔ kcmerrill (v0.2) alfred ] cd demo/
09:53 PM ✔ kcmerrill (v0.2) demo ] alfred env @themayor
[ 0s] (25 Dec 17 21:54 MST) env started [@themayor] Lets set some env variables!
[ 0s] (25 Dec 17 21:54 MST) env set $WHO kcmerrill
[ 0s] (25 Dec 17 21:54 MST) env set $TWITTER @themayor
[ 0s] (25 Dec 17 21:54 MST) twitter:@themayor
[ 0s] (25 Dec 17 21:54 MST) who:kcmerrill
[ 0s] (25 Dec 17 21:54 MST) env ✔ ok [@themayor] elapsed time '0s'
09:54 PM ✔ kcmerrill (v0.2) demo ]
```

### serve | string(port)

This component will allow you to serve static context based on `dir`. The string provided will be the port, and the server will only last for as long as tasks are running. 

By default, the server will bind to `0.0.0.0:<port>`, however, you can change this by using `127.0.0.1:<port>` where `127.0.0.1` is what you want bound.

```yaml
static.web.server:
    summary: Lets start up a static web server
    serve: 8091
    command: |
      curl --verbose http://localhost:8091/myfile.txt
```

```sh
10:12 PM ✔ kcmerrill (v0.2) demo ] echo "myfile" > myfile.txt
10:12 PM ✔ kcmerrill (v0.2) demo ] alfred static.web.server
[ 0s] (25 Dec 17 22:13 MST) static.web.server started [] Lets start up a static web server
[ 0s] (25 Dec 17 22:13 MST) static.web.server serving ./ 0.0.0.0:8091
[ 0s] (25 Dec 17 22:13 MST)   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
[ 0s] (25 Dec 17 22:13 MST)                                  Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0*   Trying ::1...
[ 0s] (25 Dec 17 22:13 MST) * TCP_NODELAY set
[ 0s] (25 Dec 17 22:13 MST) * Connected to localhost (::1) port 8091 (#0)
[ 0s] (25 Dec 17 22:13 MST) > GET /myfile.txt HTTP/1.1
[ 0s] (25 Dec 17 22:13 MST) > Host: localhost:8091
[ 0s] (25 Dec 17 22:13 MST) > User-Agent: curl/7.54.0
[ 0s] (25 Dec 17 22:13 MST) > Accept: */*
[ 0s] (25 Dec 17 22:13 MST) >
[ 0s] (25 Dec 17 22:13 MST) < HTTP/1.1 200 OK
[ 0s] (25 Dec 17 22:13 MST) < Accept-Ranges: bytes
[ 0s] (25 Dec 17 22:13 MST) myfile
[ 0s] (25 Dec 17 22:13 MST) < Content-Length: 7
[ 0s] (25 Dec 17 22:13 MST) < Content-Type: text/plain; charset=utf-8
[ 0s] (25 Dec 17 22:13 MST) < Last-Modified: Tue, 26 Dec 2017 05:12:41 GMT
[ 0s] (25 Dec 17 22:13 MST) < Date: Tue, 26 Dec 2017 05:13:52 GMT
[ 0s] (25 Dec 17 22:13 MST) <
[ 0s] (25 Dec 17 22:13 MST) { [7 bytes data]
100     7  100     7    0     0    197      0 --:--:-- --:--:-- --:--:--   200
[ 0s] (25 Dec 17 22:13 MST) * Connection #0 to host localhost left intact
[ 0s] (25 Dec 17 22:13 MST) static.web.server ✔ ok [] elapsed time '0s'
10:13 PM ✔ kcmerrill (v0.2) demo ]
```

### setup | TaskGroup{}

A [taskgroup](#taskgroups) that gets run at the start of a task. The given tasks are run in the order they are provided.

```yaml
setup.task.one:
    summary: one task
    command: echo one task

setup.task.two:
    summary: two task
    command: echo two task {{ index .Args 0 }}

setup.task:
    summary: This is the main task
    setup: |
        setup.task.one
        setup.task.two({{ index .Args 0 }})
```

```sh
10:24 PM ✘ kcmerrill (v0.2) demo ] alfred setup.task arg.one
[ 0s] (25 Dec 17 22:24 MST) setup.task started [arg.one] This is the main task
[ 0s] (25 Dec 17 22:24 MST) setup.task setup setup.task.one, setup.task.two
[ 0s] (25 Dec 17 22:24 MST) setup.task.one started [] one task
[ 0s] (25 Dec 17 22:24 MST) one task
[ 0s] (25 Dec 17 22:24 MST) setup.task.one ✔ ok [] elapsed time '0s'
[ 0s] (25 Dec 17 22:24 MST) setup.task.two started [arg.one] two task
[ 0s] (25 Dec 17 22:24 MST) two task arg.one
[ 0s] (25 Dec 17 22:24 MST) setup.task.two ✔ ok [arg.one] elapsed time '0s'
[ 0s] (25 Dec 17 22:24 MST) setup.task ✔ ok [arg.one] elapsed time '0s'
10:24 PM ✔ kcmerrill (v0.2) demo ]
```

### multitask | TaskGroup{}

A [taskgroup](#taskgroups) that gets run concurrently. 

Please note, do due it's concurrency, do not rely on arguments that maybe be updated or relied upon in other multitask tasks. Also, careful with the `dir` component in the given multitasks as you might get odd results of directories switching unexpectedly. 

```yaml
multi.task.one:
    summary: one task
    command: echo one task

multi.task.two:
    summary: two task
    command: echo two task {{ index .Args 0 }}

multi.task:
    summary: Run tasks concurrently
    multitask: |
        multi.task.one
        multi.task.two({{ index .Args 0 }})
```

```sh
03:01 PM ✘ kcmerrill (v0.2) demo ] alfred multi.task another.task
[ 0s] (26 Dec 17 15:01 MST) multi.task started [another.task] Run tasks concurrently
[ 0s] (26 Dec 17 15:01 MST) multi.task multitasks multi.task.one, multi.task.two
[ 0s] (26 Dec 17 15:01 MST) multi.task.one started [] one task
[ 0s] (26 Dec 17 15:01 MST) multi.task.two started [another.task] two task
[ 0s] (26 Dec 17 15:01 MST) two task another.task
[ 0s] (26 Dec 17 15:01 MST) multi.task.two ✔ ok [another.task] elapsed time '0s'
[ 1s] (26 Dec 17 15:01 MST) one task
[ 1s] (26 Dec 17 15:01 MST) multi.task.one ✔ ok [] elapsed time '1s'
[ 1s] (26 Dec 17 15:01 MST) multi.task ✔ ok [another.task] elapsed time '1s'
```

### tasks | TaskGroup{}

A [taskgroup](#taskgroups) that runs in order the tasks are provided.

If you are not using `multitask`, this component and `setup` are essentially the same and you can use either. 

```yaml
tasks.task.one:
    summary: one task
    command: echo one task

tasks.task.two:
    summary: two task
    command: echo two task {{ index .Args 0 }}

tasks:
    summary: Run tasks before running a command
    tasks: |
        tasks.task.one
        tasks.task.two({{ index .Args 0 }})
```

```sh
03:05 PM ✘ kcmerrill (v0.2) demo ] alfred tasks another.task
[ 0s] (26 Dec 17 15:06 MST) tasks started [another.task] Run tasks before running a command
[ 0s] (26 Dec 17 15:06 MST) tasks tasks tasks.task.one, tasks.task.two
[ 0s] (26 Dec 17 15:06 MST) tasks.task.one started [] one task
[ 0s] (26 Dec 17 15:06 MST) one task
[ 0s] (26 Dec 17 15:06 MST) tasks.task.one ✔ ok [] elapsed time '0s'
[ 0s] (26 Dec 17 15:06 MST) tasks.task.two started [another.task] two task
[ 0s] (26 Dec 17 15:06 MST) two task another.task
[ 0s] (26 Dec 17 15:06 MST) tasks.task.two ✔ ok [another.task] elapsed time '0s'
[ 0s] (26 Dec 17 15:06 MST) tasks ✔ ok [another.task] elapsed time '0s'
03:06 PM ✔ kcmerrill (v0.2) demo ]
```

### watch | string(regExp)

When set, will pause the task and look at modified times every second for the given regular expression. Given the regular expression `.*?go$` will pause and wait for any file that ends with `.go` to be modified, before continuing on. 

```yaml
tdd.go:
    summary: Watch .go files and run test.go
    watch: ".*?go$"
    command: go test $(go list ./... | grep -v /vendor/) --race
    # will only run once, set every: <duration> to repeat
    #every: 1s
```

```sh
03:20 PM ✔ kcmerrill (v0.2) alfred ] alfred tdd.go
[ 0s] (26 Dec 17 15:21 MST) tdd.go started [] Watch .go files and run test.go
[ 0s] (26 Dec 17 15:21 MST) tdd.go watching ./
[ 3s] (26 Dec 17 15:21 MST) tdd.go modified alfred/every.go
[ 6s] (26 Dec 17 15:21 MST) ?        github.com/kcmerrill/alfred     [no test files]
[10s] (26 Dec 17 15:21 MST) ok          github.com/kcmerrill/alfred/alfred      4.250s
[11s] (26 Dec 17 15:21 MST) tdd.go ✔ ok [] elapsed time '11s'
03:21 PM ✔ kcmerrill (v0.2) alfred ] alfred tdd.go
[ 0s] (26 Dec 17 15:22 MST) tdd.go started [] Watch .go files and run test.go
[ 0s] (26 Dec 17 15:22 MST) tdd.go watching ./
[ 3s] (26 Dec 17 15:22 MST) tdd.go modified alfred/command.go
[ 6s] (26 Dec 17 15:22 MST) ?        github.com/kcmerrill/alfred     [no test files]
[10s] (26 Dec 17 15:22 MST) ok          github.com/kcmerrill/alfred/alfred      4.145s
[10s] (26 Dec 17 15:22 MST) tdd.go ✔ ok [] elapsed time '10s'
[10s] (26 Dec 17 15:22 MST) tdd.go every 1s
[11s] (26 Dec 17 15:22 MST) tdd.go started [] Watch .go files and run test.go
[11s] (26 Dec 17 15:22 MST) tdd.go watching ./
```

### for | map[string]string

A simple looping component containing the following keys: `tasks`, `multitask` and `args`. 

#### tasks | TaskGroup{}

A list of tasks that will be run in order, with each new line of `args` provided as the first argument.

#### multitask | TaskGroup{}

A list of tasks that will be run concurrently, with each new line of `args` provided as the first argument.

#### args | string(command, text)

The string is evaluated as a command, and if the result is a non zero exit code will be the resulting value. 

If the string is not a valid command, it will be consumed as a regular string. 

The resulting string should contain new lines which will be passed through as the first argument to the given `task` or `multitask` provided.

```yaml
for.loop.one:
    summary: Demo our for loop with plain text
    for:
      args: |
          batman
          robin
          spiderman
      tasks: |
          for.loop.echo

for.loop.two:
    summary: Demo our for loop as a command with arguments
    for:
      args: |
        ls -R -1
      tasks: |
        for.loop.echo
```

```sh
03:44 PM ✔ kcmerrill (v0.2) demo ] alfred for.loop.one
[ 0s] (26 Dec 17 15:44 MST) for.loop.one started [] Demo our for loop
[ 0s] (26 Dec 17 15:44 MST) for.loop.echo started [batman] A simple task that echos the first argument
[ 0s] (26 Dec 17 15:44 MST) batman
[ 0s] (26 Dec 17 15:44 MST) for.loop.echo ✔ ok [batman] elapsed time '0s'
[ 0s] (26 Dec 17 15:44 MST) for.loop.echo started [robin] A simple task that echos the first argument
[ 0s] (26 Dec 17 15:44 MST) robin
[ 0s] (26 Dec 17 15:44 MST) for.loop.echo ✔ ok [robin] elapsed time '0s'
[ 0s] (26 Dec 17 15:44 MST) for.loop.echo started [spiderman] A simple task that echos the first argument
[ 0s] (26 Dec 17 15:44 MST) spiderman
[ 0s] (26 Dec 17 15:44 MST) for.loop.echo ✔ ok [spiderman] elapsed time '0s'
[ 0s] (26 Dec 17 15:44 MST) for.loop.one ✔ ok [] elapsed time '0s'
03:44 PM ✔ kcmerrill (v0.2) demo ]
```

```sh
03:44 PM ✔ kcmerrill (v0.2) demo ] alfred for.loop.two
[ 0s] (26 Dec 17 15:46 MST) for.loop.two started [] Demo our for loop as a command with arguments
[ 0s] (26 Dec 17 15:46 MST) for.loop.echo started [alfred.yml] A simple task that echos the first argument
[ 0s] (26 Dec 17 15:46 MST) alfred.yml
[ 0s] (26 Dec 17 15:46 MST) for.loop.echo ✔ ok [alfred.yml] elapsed time '0s'
[ 0s] (26 Dec 17 15:46 MST) for.loop.echo started [myfile.txt] A simple task that echos the first argument
[ 0s] (26 Dec 17 15:46 MST) myfile.txt
[ 0s] (26 Dec 17 15:46 MST) for.loop.echo ✔ ok [myfile.txt] elapsed time '0s'
[ 0s] (26 Dec 17 15:46 MST) for.loop.two ✔ ok [] elapsed time '0s'
```

### command | string

A string that will be run as a part of `bash -c <command>`. 

If the command fails, the task will fail and the flow throughout the rest of the components if required. 

Please note, failed commands do _NOT_ fail the task. You will still get a `0` exit code from alfred. If you do need to immediately fail, see [exit](#exit) component.

Also, command _CAN_ take in multiple lines but only the last line of the command is evaluated with it's exit code if it's not part of a larger command. You'll need to try a few sample out to get a feel for it. If you need _EVERY_ line evaluated as it's own command, please refer to the [commands](#commands) component.

```yaml
demo.command:
  summary: This is a single line command
  command: echo "hello world!"

demo.command.two:
  summary: This is a multi line command with a failure(notice |)
  command: |
    echo I will fail on purpose && false
    echo "Notice how I still am displayed?"
demo.command.three:
    summary: Because YAML is awesome, you can do this too!(notice >)
    command: >
        docker
        run
        -ti
        --rm
        ubuntu
```

```sh
04:24 PM ✔ kcmerrill (v0.2) demo ] alfred demo.command
[ 0s] (26 Dec 17 16:24 MST) demo.command started [] This is a single line command
[ 0s] (26 Dec 17 16:24 MST) hello world!
[ 0s] (26 Dec 17 16:24 MST) demo.command ✔ ok [] elapsed time '0s'
04:24 PM ✔ kcmerrill (v0.2) demo ] alfred demo.command.two
[ 0s] (26 Dec 17 16:25 MST) demo.command.two started [] This is a multi line command with a failure
[ 0s] (26 Dec 17 16:25 MST) I will fail on purpose
[ 0s] (26 Dec 17 16:25 MST) Notice how I still am displayed?
[ 0s] (26 Dec 17 16:25 MST) demo.command.two ✔ ok [] elapsed time '0s'
04:25 PM ✔ kcmerrill (v0.2) demo ]
```

### commands | string

A new line separated string that will be run as a part of `bash -c <command>`. Identical to command, however, use commands component if you need every command to be evaulated for pass/fail. 

### ok | TaskGroup{}

A [taskgroup](#taskgroups) that runs in order the tasks are provided only if the task is ok up to this point(meaning that everything has been OK)

```yaml
ok.task.one:
    summary: one task
    command: echo one task

ok.task.two:
    summary: two task
    command: echo two task {{ index .Args 0 }}

ok.tasks:
    summary: run tasks before running a command
    command: |
        echo hello world
    ok: |
        tasks.task.one
        tasks.task.two({{ index .Args 0 }})
```

```sh
06:57 PM ✘ kcmerrill (v0.2) demo ] alfred ok.tasks second.ok.task
[ 0s] (26 Dec 17 18:58 MST) ok.tasks started [second.ok.task] run tasks before running a command
[ 0s] (26 Dec 17 18:58 MST) hello world
[ 0s] (26 Dec 17 18:58 MST) ok.tasks ✔ ok [second.ok.task] elapsed time '0s'
[ 0s] (26 Dec 17 18:58 MST) ok.tasks ok.tasks ok.task.one, ok.task.two
[ 0s] (26 Dec 17 18:58 MST) ok.task.one started [] one task
[ 0s] (26 Dec 17 18:58 MST) one task
[ 0s] (26 Dec 17 18:58 MST) ok.task.one ✔ ok [] elapsed time '0s'
[ 0s] (26 Dec 17 18:58 MST) ok.task.two started [second.ok.task] two task
[ 0s] (26 Dec 17 18:58 MST) two task second.ok.task
[ 0s] (26 Dec 17 18:58 MST) ok.task.two ✔ ok [second.ok.task] elapsed time '0s'
06:58 PM ✔ kcmerrill (v0.2) demo ]
```

### http.tasks | map[string]string

Baked in is the ability to run tasks via HTTP requests. Specify the `port` and if necessary `password`. The password field will be translated and evaluated, so feel free to use env vars, scripts etc. The output will be sent to the http response instead of stdout. If you want to see what is happening, you can always log the output as well. 

```yaml
http.task:
    summary: HTTP Task runner
    http.tasks:
        port: 8093
        password: password
```

```sh
01:15 AM ✔ kcmerrill (master) alfred ] alfred http.task
[    0s] (28 Dec 17 01:15 MST) http.task started [] HTTP Task runner
[    0s] (28 Dec 17 01:15 MST) http.task serving ./ 0.0.0.0:8093
[    4s] (28 Dec 17 01:15 MST) http.task build started [one, two, three]
```

```sh
01:21 AM ✘ kcmerrill  tmp ] curl -u "password:" -X POST -d "body.of.msg" http://localhost:8093/echo/alfredwashere
[ 1m40s] (28 Dec 17 01:21 MST) echo started [alfredwashere] echo out some shit
[ 1m40s] (28 Dec 17 01:21 MST) args: alfredwashere
[ 1m40s] (28 Dec 17 01:21 MST) stdin: body.of.msg
[ 1m40s] (28 Dec 17 01:21 MST) echo  ok [alfredwashere] elapsed time '1m40s'
01:21 AM ✔ kcmerrill  tmp ] curl -X POST -d "body.of.msg" http://localhost:8093/echo/alfredwashere
{"error": "unauthorized"}
01:21 AM ✔ kcmerrill  tmp ]
```

### fail | TaskGroup{}

A [taskgroup](#taskgroups) that runs in order the tasks are provided only if the task has failed up to this point.

```yaml
fail.task.one:
    summary: one task
    command: echo one task

fail.task.two:
    summary: two task
    command: echo two task {{ index .Args 0 }}

fail.tasks:
    summary: Demonstrate failed tasks
    command: |
        echo hello world
    fail: |
        fail.task.one
        fail.task.two({{ index .Args 0 }})
```

```sh
07:00 PM ✔ kcmerrill (v0.2) demo ] alfred fail.tasks second.fail.task
[ 0s] (26 Dec 17 19:00 MST) fail.tasks started [second.fail.task] Demonstrate failed tasks
[ 0s] (26 Dec 17 19:00 MST) ls: /tmp/idonotexist: No such file or directory
[ 0s] (26 Dec 17 19:00 MST) fail.tasks command failed
[ 0s] (26 Dec 17 19:00 MST) fail.tasks ✘ failed [second.fail.task] elapsed time '0s'
[ 0s] (26 Dec 17 19:00 MST) fail.tasks fail.tasks fail.task.one, fail.task.two
[ 0s] (26 Dec 17 19:00 MST) fail.task.one started [] one task
[ 0s] (26 Dec 17 19:00 MST) one task
[ 0s] (26 Dec 17 19:00 MST) fail.task.one ✔ ok [] elapsed time '0s'
[ 0s] (26 Dec 17 19:00 MST) fail.task.two started [second.fail.task] two task
[ 0s] (26 Dec 17 19:00 MST) two task second.fail.task
[ 0s] (26 Dec 17 19:00 MST) fail.task.two ✔ ok [second.fail.task] elapsed time '0s'
07:00 PM ✔ kcmerrill (v0.2) demo ]
```

### wait | duration

Once a task has completed, exits, before running again, you can wait a [golang duration](https://golang.org/pkg/time/#ParseDuration) before the next task runs. 

```yaml
07:03 PM ✘ kcmerrill (v0.2) demo ] alfred wait
[ 0s] (26 Dec 17 19:07 MST) wait started [] Wait a golang duration before continuing
[ 0s] (26 Dec 17 19:07 MST) waiting
[ 0s] (26 Dec 17 19:07 MST) wait ✔ ok [] elapsed time '0s'
[ 0s] (26 Dec 17 19:07 MST) wait wait 10s
07:07 PM ✔ kcmerrill (v0.2) demo ]
```

```sh
07:03 PM ✘ kcmerrill (v0.2) demo ] alfred wait
[ 0s] (26 Dec 17 19:07 MST) wait started [] Wait a golang duration before continuing
[ 0s] (26 Dec 17 19:07 MST) waiting
[ 0s] (26 Dec 17 19:07 MST) wait ✔ ok [] elapsed time '0s'
[ 0s] (26 Dec 17 19:07 MST) wait wait 10s
07:07 PM ✔ kcmerrill (v0.2) demo ]
```

### every | duration

A lot like the [wait](#wait) component, but instead of exiting, will rerun the same task again.


```yaml
every:
  summary: Run a command every <golang duration>
  every: "{{ index .Args 0 }}"
  command: |
    echo every {{ index .Args 0 }}
```

```sh
07:07 PM ✔ kcmerrill (v0.2) demo ] alfred every 1s
[ 0s] (26 Dec 17 19:11 MST) every started [1s] Run a command every <golang duration>
[ 0s] (26 Dec 17 19:11 MST) every 1s
[ 0s] (26 Dec 17 19:11 MST) every ✔ ok [1s] elapsed time '0s'
[ 0s] (26 Dec 17 19:11 MST) every every 1s
[ 1s] (26 Dec 17 19:11 MST) every started [1s] Run a command every <golang duration>
[ 1s] (26 Dec 17 19:11 MST) every 1s
[ 1s] (26 Dec 17 19:11 MST) every ✔ ok [1s] elapsed time '1s'
[ 1s] (26 Dec 17 19:11 MST) every every 1s
[ 2s] (26 Dec 17 19:11 MST) every started [1s] Run a command every <golang duration>
[ 2s] (26 Dec 17 19:11 MST) every 1s
[ 2s] (26 Dec 17 19:11 MST) every ✔ ok [1s] elapsed time '2s'
[ 2s] (26 Dec 17 19:11 MST) every every 1s
^C
07:11 PM ✘ kcmerrill (v0.2) demo ]
```

# Tips and Tricks

Here are some tips and tricks that are regularly used, or things that have replaced functionality in previous versions of Alfred.

### Aliases

Lets say you have a task and want it named a few different ways. There are a ton of reasons why you might want to do that. In previous versions of alfred, you'd use an `alias` component, but it added extra complexity that wasn't needed. You can do this quite simply with [YAML Anchors](https://en.wikipedia.org/wiki/YAML#Advanced_components).

```yaml
hello.world:
    summary: I will say hello world!
    command: |
        echo hello world!

task.original: &task_original
    summary: Original task
    setup: hello.world
    command: |
        echo I am the original task.

task.alias: *task_original
```

### Task Inheritance 

Lets say you've created a rather large and complex task. Sure, you could copy it, but if you change one you might need to change them both. Not ideal solution. You can use [YAML Anchors](https://en.wikipedia.org/wiki/YAML#Advanced_components).

In the example below, any `key` below `<<: *task_original` will inherit from the original task, and `command` component will override the `task_original` command. 

As a side note, YAML is not a huge fan of dots(`.`) in anchor names which is why I used underscores(`_`) in the anchor name.

```yaml
hello.world:
    summary: I will say hello world!
    command: |
        echo hello world!

task.original: &task_original
    summary: Original task
    setup: hello.world
    command: |
        echo I am the original task.

task.alias:
    <<: *task_original
    command: |
        echo "I am the alias"
```

### Run a task X times

If you've taken a peek at the `for` component, you'll find that it allows tasks to iterate over data. Sometimes you just need to do something `x` times. Load testing, batching etc. To do that, you can use the for loop built in, along with golang templates to create a numerical loop by supplying a `range` in the `for` component, like so: `{{range $i, $e := until 5}}{{$i}}\n{{end}}`. You can read more about it at [masterminds/sprig](https://github.com/Masterminds/sprig).

```yaml
echo:
    command: echo {{ index .Args 0 }}
range:
    summary: Iterate X number of times
    for:
        args: "{{range $i, $e := until 5}}{{$i}}\n{{end}}"
        tasks: echo
```
