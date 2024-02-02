
# Helm 3

Installe o helm 3 com os seguintes comandos.

```bash
sudo su

HELM_TAR_FILE=helm-v3.12.0-linux-amd64.tar.gz
HELM_URL=https://get.helm.sh
HELM_BIN=helm3

function install_helm3 {

if [ -z $(which $HELM_BIN) ]; then
    wget ${HELM_URL}/${HELM_TAR_FILE}
    tar -xvzf ${HELM_TAR_FILE}
    chmod +x linux-amd64/helm
    sudo cp linux-amd64/helm /usr/local/bin/$HELM_BIN
    sudo ln -sfn /usr/local/bin/$HELM_BIN /usr/local/bin/helm
    rm -rf ${HELM_TAR_FILE} linux-amd64
    echo -e "\nwhich ${HELM_BIN}"
    which ${HELM_BIN}
else
    echo "Helm 3 is most likely installed"
fi
}

install_helm3

which $HELM_BIN

$HELM_BIN version

exit
```

Referências:
* https://helm.sh/docs/