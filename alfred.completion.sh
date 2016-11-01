# Bash completion for alfred

_alfred()
{
    local cur prev tasks flags
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    tasks="$(alfred | grep '^\[' | cut -d ' ' -f 1 | tr -d '[]' | tr '\n' ' ')"
    flags="-dir -port -serve"

    case "${cur}" in
        -*)
            # looks like a -flag
            COMPREPLY=( $(compgen -W "$flags" -- ${cur}) )
            return 0
            ;;
        *)
            if echo " $tasks " | grep -q " $cur"; then
                # current word is a prefix of a valid task
                COMPREPLY=( $(compgen -W "$tasks" -- ${cur}) )
                return 0
            else
                # fall back to default (path) completion
                return 1
            fi
            ;;
    esac
}

complete -o default -F _alfred alfred
