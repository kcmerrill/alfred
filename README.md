# alfred
Because even Batman needs a little help.
![Alfred](https://raw.githubusercontent.com/kcmerrill/alfred/master/alfred.jpg "Alfred")

A simple go/yaml powered make file with a bit of a twist. 

* Tasks can call other tasks
* Tasks can be run every <insert duration>
* Tasks can pause
* Tasks can call other tasks depending on success/failure
* Alfred uses go templates so you can inject variables making tasks reusable
* Common tasks can be invoked inside using alfred
* Common tasks will be stored in this repository for shared use(git, docker, mail, slack, etc ...) (coming soon)
* Optional private/public repositories so you can share private tasks with coworkers


## Why
I have a lot of tasks I do daily that I'd like a remote repository to use. Not only can I update them, but others can also use it too. I love make files, but wanted just a bit more functionality while keeping it simple. 

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
