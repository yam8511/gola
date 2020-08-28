# # kind
# reg_name='kind-gola-registry'
# kind delete cluster --name gola || exit 0
# docker rm -f ${reg_name} || exit 0
# docker network rm kind || exit 0

# k3d
reg_name='k3d-gola-registry'
docker rm -f ${reg_name} || exit 0
k3d cluster delete gola
docker network rm k3d-gola || exit 0
