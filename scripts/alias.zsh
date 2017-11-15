# List and select pod name with fzf (https://github.com/junegunn/fzf)
# e.g.
#   kubectl exec -it P sh
#   kubectl delete pod P
alias fzfkubernetesalias="fzf --height 25 --header-lines=1 --reverse --multi --cycle"
alias -g P='$(kubectl get pods | fzfkubernetesalias | awk "{print \$1}")'

# Like P, global aliases about kubernetes resources
alias -g PO='$(kubectl get pods | fzfkubernetesalias | awk "{print \$1}")'
alias -g NS='$(kubectl get ns   | fzfkubernetesalias | awk "{print \$1}")'

# alias -g DEPLOY='$(kubectl get deploy| fzf-tmux --header-lines=1 --reverse --multi --cycle | awk "{print \$1}")'
# alias -g     RS='$(kubectl get rs    | fzf-tmux --header-lines=1 --reverse --multi --cycle | awk "{print \$1}")'
# alias -g    SVC='$(kubectl get svc   | fzf-tmux --header-lines=1 --reverse --multi --cycle | awk "{print \$1}")'
# alias -g    ING='$(kubectl get ing   | fzf-tmux --header-lines=1 --reverse --multi --cycle | awk "{print \$1}")'

# Context switcher
# c.f. https://github.com/ahmetb/kubectx
alias kubectl-change='kubectx $(kubectx | fzy) >/dev/null'

# References
# - https://github.com/c-bata/kube-prompt
# - https://github.com/cloudnativelabs/kube-shell
