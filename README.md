kubetools
=========

A toolkit for k8s developers

## Usage

### To begin with

Some convenient commands are in `bin` directory. So you should add that directory to `$PATH` at first.

```bash
export PATH="$PWD/bin:$PATH"
```

### For CLI

- kube
- kubectx (extended of [ahmetb/kubectx](https://github.com/ahmetb/kubectx))
- kubens (ditto)

### For tmux

- kube-context
- gcp-context

```config
set-option -g status-left 'tmux:[#P] #[fg=colour33](K) #(kube-context)#[default] #[fg=colour1](G) #(gcp-context)#[default]'
```

## License

MIT

## Auther

b4b4r07
