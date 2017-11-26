# List and select pod name with fzf (https://github.com/junegunn/fzf)
# e.g.
#   kubectl exec -it P sh
#   kubectl delete pod P
alias fzfkubernetesalias="fzf --height 25 --header-lines=1 --reverse --multi --cycle"
alias -g P='$(kubectl get pods | fzfkubernetesalias | awk "{print \$1}")'

# Like P, global aliases about kubernetes resources
alias -g  PO='$(kubectl get pods | fzfkubernetesalias | awk "{print \$1}")'
alias -g  NS='$(kubectl get ns   | fzfkubernetesalias | awk "{print \$1}")'
alias -g  RS='$(kubectl get rs   | fzfkubernetesalias | awk "{print \$1}")'
alias -g SVC='$(kubectl get svc  | fzfkubernetesalias | awk "{print \$1}")'
alias -g ING='$(kubectl get ing  | fzfkubernetesalias | awk "{print \$1}")'

# References
# - https://github.com/c-bata/kube-prompt
# - https://github.com/cloudnativelabs/kube-shell
