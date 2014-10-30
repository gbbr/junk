##
# MAMP SETUP AND COLORS
##

export GOPATH=~/go
export PATH=$GOPATH/bin:$PATH
export CLICOLOR=1
export TERM=xterm-256color
export CDPATH=".:$GOPATH/src/github.com:$CDPATH"
export EDITOR=vim

##
# COMMAND LINE FORMAT
##

function parse_git_branch () {
       git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/ (\1)/'
}

if [ -f ~/.git-completion.bash ]; then
  . ~/.git-completion.bash
fi
 
RED="\[\033[0;31m\]"
YELLOW="\[\033[0;33m\]"
GREEN="\[\033[0;32m\]"
NO_COLOUR="\[\033[0m\]"
 
PS1="$GREEN\u@\h$NO_COLOUR:\W$YELLOW\$(parse_git_branch)$NO_COLOUR\$ "

##
# MY ALIASES
##

alias bye='exit'
alias swap='echo; echo "Breakdown:"; du -hs /private/var/vm/*; echo; echo "Total: "; du -hs /private/var/vm/;echo;'
alias gg="git grep -C 3 $@"

##
# HISTORY AUCO-COMPLETE
##

bind '"\e[A": history-search-backward'
bind '"\e[B": history-search-forward'

export NVM_DIR="/Users/gabriel/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"  # This loads nvm
