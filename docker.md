# Docker

Instale o Docker CE (Community Edition) com os seguintes comandos:

```bash
curl -fsSL https://get.docker.com | sudo bash

# References:
# https://stackoverflow.com/questions/18836853/sudo-cat-eof-file-doesnt-work-sudo-su-does
# https://earthly.dev/blog/managing-k8s-with-kubeadm/
sudo bash -c 'cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF'


# Add your user to the Docker group
sudo apt install -y acl uidmap
dockerd-rootless-setuptool.sh install
sudo usermod -aG docker $USER
sudo setfacl -m user:$USER:rw /var/run/docker.sock
sudo setfacl -m user:$USER:rw /var/run/containerd/containerd.sock

# Start the Docker service
sudo systemctl start docker
sudo systemctl start containerd

# Configure Docker to boot up with the OS
sudo systemctl enable docker
sudo systemctl enable containerd
```
