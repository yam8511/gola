k3d cluster create gola \
    --servers 1 \
    --agents 3 \
    --api-port 0.0.0.0:6550 \
    -v $(PWD)/_setup/kube/k3d/registry.yaml:/etc/rancher/k3s/registries.yaml \
    -p 80:80@loadbalancer \
    -p 443:443@loadbalancer \
    -p 30007-30008:30007-30008@server[0]

# create registry container unless it already exists
reg_name='k3d-gola-registry'
reg_port='5000'
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
    docker run \
    -d --restart=always \
    -p "${reg_port}:5000" \
    --name "${reg_name}" \
    registry:2

    docker network connect "k3d-gola" "${reg_name}"
fi
