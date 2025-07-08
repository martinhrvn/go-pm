#compdef gopm

# Zsh completion for gopm
# Copy to your zsh completion directory or source this file

_gopm() {
    local context state line
    
    _arguments \
        '1: :->command' \
        '*: :->args'
    
    case $state in
        command)
            _values 'gopm commands' \
                'run[Interactive command selection and execution]' \
                'list[List all available commands]' \
                'help[Show help message]'
            ;;
    esac
}

_gopm "$@"