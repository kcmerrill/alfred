# The Manual

Alfred tasks can be complex or as simple as you want them to be. The idea is they are small, reusable and can help automate your daily workflow. Each task is made up of [Components](#components). You can think of these components as building blocks to build whatever it is you need. In this manual, we'll go over some classic use cases for alfred, we'll discuss each component in depth, tips and tricks and also ways of structuring your tasks for maximum awesomeness.

## Components 

The components here will be listed in order in which they are executed within Alfred. With the way golang's maps work, [they are randomized](https://github.com/golang/go/issues/2630) to prevent DOS attacks. The reason this is important is, your yaml file and it's tasks can be ordered however you'd like, but they will be executed in a specific order. 

### log | string

The raw output of commands will be sent to the specified filename appened, created if the file does not exist.

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

### dir | string

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

### config | string

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

### register | map[string]string 

Based on `key: value` pairs, will register the pairs as [variables](#variables). The value is then `evaluated` when a zero exit code is shown, the `CombinedOutput()` is the resulting value. 

```yaml
09:43 PM ✔ kcmerrill (v0.2) demo ] alfred register
[ 0s] (25 Dec 17 21:43 MST) register started [] Demonstrate the registration of variables
[ 0s] (25 Dec 17 21:43 MST) register registered user kcmerrill
[ 0s] (25 Dec 17 21:43 MST) register registered twitter @themayor
[ 0s] (25 Dec 17 21:43 MST) kcmerrill
[ 0s] (25 Dec 17 21:43 MST) @themayor
[ 0s] (25 Dec 17 21:43 MST) register ✔ ok [] elapsed time '0s'
09:43 PM ✔ kcmerrill (v0.2) demo ]
```

```sh
register:
    summary: Demonstrate the registration of variables
    register:
        user: whoami
        twitter: "@themayor"
    command: |
      echo "{{ index .Vars "user" }}"
      echo "{{ .Vars.twitter }}"
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

### serve | string 

This component will allow you to serve static context based on `dir`. The string provided will be the port, and the server will only last for as long as tasks are running. 

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

A [taskgroup](#taskgroup) that gets run at the start of a task. 

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