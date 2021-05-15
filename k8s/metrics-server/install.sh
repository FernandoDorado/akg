helm repo add bitnami https://charts.bitnami.com/bitnami
helm upgrade \
    -i \
    -n kube-system \
    metrics-server \
    bitnami/metrics-server \
    -f values.yml
