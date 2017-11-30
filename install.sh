#!/bin/bash

set -ex

install_dir=~/bin

install_kubectx() {
    cmd=kubectx
    download_file=/tmp/$cmd
    {
        echo "#!/bin/bash"
        curl "https://raw.githubusercontent.com/ahmetb/kubectx/master/utils.bash"
        curl "https://raw.githubusercontent.com/ahmetb/kubectx/master/$cmd" | sed -e 's/source/: source/g'
    } >$download_file 2>/dev/null
    chmod 755 $download_file
    install -m 755 $download_file $install_dir/kube--ctx
    install -m 755 $cmd/$cmd $install_dir
}

install_kubens() {
    cmd=kubens
    download_file=/tmp/$cmd
    {
        echo "#!/bin/bash"
        curl "https://raw.githubusercontent.com/ahmetb/kubectx/master/utils.bash"
        curl "https://raw.githubusercontent.com/ahmetb/kubectx/master/$cmd" | sed -e 's/source/: source/g'
    } >$download_file 2>/dev/null
    chmod 755 $download_file
    install -m 755 $download_file $install_dir/kube--ns
    install -m 755 $cmd/$cmd $install_dir
}

case $1 in
    "kubectx")
        install_kubectx
        ;;
    "kubens")
        install_kubens
        ;;
    "")
        install_kubectx
        install_kubens
        ;;
    *)
        echo "nothing"
        ;;
esac
