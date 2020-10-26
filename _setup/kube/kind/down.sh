echo "⚙️  請輸入 kubernetes cluster 名稱 [預設: default]"
printf "> "
read cluster
cluster=${cluster:-default}

kind delete cluster --name ${cluster} || exit 0

K8S_CLUSTER_FD=$HOME/.cache/gola/k8s/cluster-${cluster}
reg_name='registry'
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"

if [ "${running}" = 'true' ]; then
    printf "需要刪除 %s 嗎？(Y/[n])\n> " ${reg_name}
    read yes
    yes=${yes:-n}
    if [ $yes = "Y" ] || [ $yes = "y" ]; then
        docker rm -f ${reg_name}
    else
        docker network disconnect kind registry 
    fi
    docker network rm kind
fi

if ls ${K8S_CLUSTER_FD} 1>/dev/null 2>&1; then
    printf "需要刪除快取資料 %s 嗎？(Y/[n])\n> " ${K8S_CLUSTER_FD}
    read yes
    yes=${yes:-n}
    if [ $yes = "Y" ] || [ $yes = "y" ]; then
        rm -r ${K8S_CLUSTER_FD} || sudo rm -r ${K8S_CLUSTER_FD}
    fi
fi
