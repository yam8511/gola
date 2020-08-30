#! /bin/sh
# 文字樣式：
# 0 一般樣式    # 1 高亮度  # 4 底線    # 5 灰底    # 7 文字 與 背景 顏色對調

# 文字顏色：
# 30 黑色    # 31 紅色    # 32 綠色    # 33 黃色
# 34 藍色    # 35 紫色    # 36 青綠色  # 37 白色

# 背景顏色：
# 40 黑色    # 41 紅色    # 42 綠色    # 43 黃色
# 44 藍色    # 45 紫色    # 46 青綠色  # 47 白色

CRE='\033[0m'
BGWHITE='\033[47m'
CBLACK='\033[0;30m'
CGREEN='\033[1;32m'
CTEAL='\033[1;36m'
CBLUE='\033[1;34m'
CPURPLE='\033[1;35m'
CYELLOW='\033[1;33m'
BPURPLE='\033[1;45m'
BTEAL='\033[1;44m'
BBLUE='\033[1;46m'

tool=$1
opt=$2
ons=$3

cmdOpts=(
    ["1"]="up"
    ["2"]="down"
    ["3"]="addons"
    ["4"]="install"
    ["5"]="exit"
)

cmdTxt=(
    ["1"]="🚀 啟動本地kube cluster"
    ["2"]="⛔️ 關閉本地kube cluster"
    ["3"]="🧩 安裝kube擴充插件"
    ["4"]="🤖 安裝kube工具"
    ["5"]="👋 結束腳本"
)

kubeTool=(
    ["1"]="k3d"
    ["2"]="kind"
    ["3"]="exit"
)

kubeToolTxt=(
    ["1"]="🛳  k3d  (light k8s in docker)"
    ["2"]="🚢 kind (official k8s in docker)"
    ["3"]="上一頁"
)

kubeToolExe=(
    ["1"]="echo 🕹  請輸入以下指令，安裝${CTEAL}k3d${CRE}\ncurl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash"
    ["2"]="echo 🤖 ${CPURPLE}On Linux:${CRE}\n${BPURPLE}\ncurl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.8.1/kind-linux-amd64\nchmod +x ./kind\nmv ./kind /usr/local/bin/kind\n${CRE}\n\n🤖 ${CTEAL}On Mac:${CRE}\n${BTEAL}\ncurl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.8.1/kind-darwin-amd64\nchmod +x ./kind\nmv ./kind /usr/local/bin/kind\n${CRE}\n\n🤖 ${CBLUE}On Windows:${CRE}\n${CBLACK}${BBLUE}\ncurl.exe -Lo kind-windows-amd64.exe https://kind.sigs.k8s.io/dl/v0.8.1/kind-windows-amd64\nMove-Item .\kind-windows-amd64.exe c:\\\%PATH%\kind.exe\n${CRE}\n${CYELLOW}ps. %PATH% 要手動改成系統可以判斷到的資料夾${CRE}\n\n\n詳細安裝請見 => https://kind.sigs.k8s.io/docs/user/quick-start/"
)

kubeUp=(
    ["1"]="sh ./_setup/kube/k3d/up.sh"
    ["2"]="sh ./_setup/kube/kind/up.sh"
)

kubeDown=(
    ["1"]="sh ./_setup/kube/k3d/down.sh"
    ["2"]="sh ./_setup/kube/kind/down.sh"
)

addonsOpt=(
    ["1"]="kube-dashboard"
    ["2"]="kube-state-metrics"
    ["3"]="istio"
    ["4"]="rancher"
    ["5"]="exit"
)

addonsUp=(
    ["1"]="sh ./_setup/addons/dashboard/up.sh"
    ["2"]="sh ./_setup/addons/kube-state-metrics/up.sh"
    ["3"]="sh ./_setup/addons/istio/up.sh"
    ["4"]="sh ./_setup/addons/rancher/up.sh"
)

addonsDown=(
    ["1"]="sh ./_setup/addons/dashboard/down.sh"
    ["2"]="sh ./_setup/addons/kube-state-metrics/down.sh"
    ["3"]="sh ./_setup/addons/istio/down.sh"
    ["4"]="sh ./_setup/addons/rancher/down.sh"
)

addonsOptExe=(
    ["4"]="echo 請輸入以下指令，安裝${CTEAL}Istio${CRE}\n> curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.7.0 sh -\n> mv istio-1.7.0 \$HOME\n\n----\n加入以下環境變數到 $HOME/.bashrc 或 $HOME/.zshrc\n\n# Istio\nexport ISTIO_PATH=\$HOME/istio-1.7.0\nexport PATH=\$PATH:\$ISTIO_PATH/bin\nsource \$ISTIO_PATH/tools/_istioctl"
)

# wait user any key...
waitKey() {
    echo "\n💬 ...按任意鍵繼續..."
    read
}

# exe is not found.
exeNotFind() {
    if which $1  1>/dev/null 2>&1 
    then
        return 1
    else
        return 0
    fi
}

# 判斷陣列arr是否包含index
isset() {
    arr=$1
    index=$2
    if [[ ${arr[$index]} = "" ]]
    then
        return 1
    else
        return 0
    fi  
}

# 顯示使用說明
usage() {
    echo "Kubernetes For Local Cluster"
    for key in ${!cmdOpts[@]};do
        printf "${CGREEN}$key) %-10s${CRE}" ${cmdOpts[$key]}
        echo ${cmdTxt[$key]}
    done
    printf "> "
    read tool
}

waitKube(){
    for i in ${!kubeTool[@]}; do
        echo "${CBLUE}$i) ${kubeToolTxt[$i]}${CRE}"
    done
    printf "> "
    read opt

    if [ "$opt" = "" ]; then
        return
    fi

    if ! [ ${kubeTool[$opt]} = "exit" ]; then
        if exeNotFind ${kubeTool[$opt]}; then
            clear
            ${kubeToolExe[$opt]}
            opt=''
            waitKey
        fi
    fi
}

waitAddons(){
    while :
    do
        clear
        echo ${cmdTxt[$tool]}
        if isset $addonsOpt $opt; then
            break
        fi

        for i in ${!addonsOpt[@]}; do
            echo "${CBLUE}$i) ${addonsOpt[$i]}${CRE}"
        done
        printf "> "
        read opt
    done


    if ! [ ${addonsOpt[$opt]} = "exit" ]; then
        if [ ${addonsOpt[$opt]} = "istio" ]; then
            if exeNotFind "istioctl"; then
                ${addonsOptExe[$opt]}
                waitKey
                return
            fi
        fi

        echo "${CTEAL}${addonsOpt[$opt]}${CRE}"
        echo "1) up"
        echo "2) down"
        echo "*) exit"
        printf "> "
        read ons
        ons=${ons:-3}
    fi
}

install() {
    while :
    do
    clear
    if isset $kubeTool $opt
    then
        case ${kubeTool[$opt]} in
        exit)
            return 0
        ;;
        *)
            if ! ${kubeToolExe[$opt]}
            then
                exit 1
            fi
            waitKey
            return 0
        ;;
        esac
    else
        echo ${cmdTxt[$tool]}
        waitKube
    fi
    done
}

up() {
    while :
    do
    clear
    if isset $kubeTool $opt
    then
        case ${kubeTool[$opt]} in
        exit)
            return 0
        ;;
        *)
            if ! ${kubeUp[$opt]}
            then
                exit 1
            fi
            waitKey
            return 0
        ;;
        esac
    else
        echo ${cmdTxt[$tool]}
        waitKube
    fi
    done
}

down() {
    while :
    do
    clear
    if isset $kubeTool $opt
    then
        case ${kubeTool[$opt]} in
        exit)
            return 0
        ;;
        *)
            if ! ${kubeDown[$opt]}
            then
                exit 1
            fi
            waitKey
            return 0
        ;;
        esac
    else
        echo ${cmdTxt[$tool]}
        waitKube
    fi
    done
}

addons() {
    while :
    do
    clear
    if isset $addonsOpt $opt
    then
        if [ ${addonsOpt[$opt]} = "exit" ]; then
            return 0
        fi

        case $ons in
        1)  # up
            if ! ${addonsUp[$opt]}; then
                exit 1
            fi
            waitKey
        ;;
        2)  # down
            if ! ${addonsDown[$opt]}; then
                exit 1
            fi
            waitKey
        ;;
        esac
        ons=''
        opt=''
    else
        echo ${cmdTxt[$tool]}
        waitAddons
        echo "what $opt"
    fi
    done
}

while :
do
    clear
    if isset $cmdOpts $tool
    then
        if ! ${cmdOpts[$tool]}
        then
            exit 1
        fi
        tool=''
        opt=''
    else
        usage
    fi
done
