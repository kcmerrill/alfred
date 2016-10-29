# Bash completion for alfred

_alfred()
{
    local cur prev tokens
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    tokens="$(alfred | grep '^\[' | cut -d ' ' -f 1 | tr -d '[]')"

    COMPREPLY=( $(compgen -W "${tokens}" -- ${cur}) )
}
complete -F _alfred alfred
