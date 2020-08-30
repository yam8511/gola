k3d cluster delete gola

reg_name='registry'
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
if [ "${running}" = 'true' ]; then
    printf "需要刪除 %s 嗎？(Y/[n])\n> " reg_name
    read yes
    yes=${yes:-n}
    if [ $yes = "Y" ] || [ $yes = "y" ]; then
        docker rm -f ${reg_name} || exit 0
        docker network rm k3d-gola || exit 0
    fi
fi
