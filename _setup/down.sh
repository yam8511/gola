reg_name='kind-gola-registry'
kind delete cluster --name gola || exit 0
docker rm -f ${reg_name} || exit 0
docker network rm kind || exit 0
