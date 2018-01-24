# Alfred v0.1 -> Alfred v0.2

There have been some changes to the latest version of alfred that I think were necessary to the code base, along with a better ability to test out components. On top of that, there were changes made to how alfred files are structured that lead to readability from within the alfred tasks. 

Here are a few of those changes. 

## Components

1. `Skip` the skip component was deprecated. The reasoning behind this was to avoid silent failures which was a common occurance. More times than not, tasks should never have stopped executing and skipped. Of course, if you feel this to be incorrect, please let me know!

1. `For` was changed, instead of using the command, it now accepts an `args` key 

## Taskgroups

This is a new concept, that is really similiar to the old way of doing task groups, however this way allows for cleaner readability and the ability to mix and match. 

Previously, any task groups, such as `setup`, `ok`, `fail`, `tasks` you could setup your task as follows: `tasks: taskone tasktwo taskthree taskfour`. This still exists! However, if you wanted to pass in arguments, lets use `taskthree` as an example, it would have to look like this: `tasks: taskone() tasktwo() taskthree(arg1) taskfour()`. Notice how, even though one task needed arguments, it required all tasks to have `()`. 

The latest version of alfred still allows non arguments to be space separated, so the very first example will still work as you expect. However, if you need to use arguments, you now need to put your task groups on newlines. 

So these are now valid:
```yaml
taskname:
   tasks: |
        taskone
        tasktwo
        taskthree(arg1)
        taskfour
```

```yaml
taskname:
    tasks: taskthree(arg1) #notice, only have 1 task w/arg? no problem!
```