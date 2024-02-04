# Optional

Add aliases in files: ``~/.bashrc`` and ``/root/.bashrc``

```bash
alias egrep='egrep --color=auto'
alias fgrep='fgrep --color=auto'
alias grep='grep --color=auto'
alias k='kubectl'
alias l='ls -CF'
alias la='ls -A'
alias live='curl parrot.live'
alias ll='ls -alF'
alias ls='ls --color=auto'
alias nettools='kubectl run --rm -it nettools-$(< /dev/urandom tr -dc a-z-0-9 | head -c${1:-4}) --image=aeciopires/nettools:1.0.0 -n default -- bash'
alias randompass='< /dev/urandom tr -dc _A-Z-a-z-0-9 | head -c${1:-16}'
alias sc='source ~/.bashrc'
alias show-hidden-files='du -sch .[!.]* * |sort -h'
alias gitlog='git log -p'
```

Apply new aliases:

```bash
source ~/.bashrc
source /root/.bashrc
```
