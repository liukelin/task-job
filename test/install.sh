#!/bin/bash  
# 测试用的脚本
underline() { printf "${underline}${bold}%s${reset}\n" "$@"
}
h1() { printf "\n${underline}${bold}${blue}%s${reset}\n" "$@"
}
h2() { printf "\n${underline}${bold}${white}%s${reset}\n" "$@"
}
debug() { printf "${white}%s${reset}\n" "$@"
}
info() { printf "${white}➜ %s${reset}\n" "$@"
}
success() { printf "${green}✔ %s${reset}\n" "$@"
}
error() { printf "${red}✖ %s${reset}\n" "$@"
}
warn() { printf "${tan}➜ %s${reset}\n" "$@"
}
bold() { printf "${bold}%s${reset}\n" "$@"
}
note() { printf "\n${underline}${bold}${blue}Note:${reset} ${blue}%s${reset}\n" "$@"
}

note "开始调用测试install脚本。"

# 分支名称
branch="${1:-$branch}"

if [ ! -n "$branch" ]; then
error "branch分支信息必须配置。"
exit 0
fi

# 如果是传参的方式，全部使用默认配置
if [ ! -n "$1" ]; then
# 部署路径
data_dir="${DATA_DIR:-/opt/operation_parcel/code/$branch}"
# 容器名
PRODUCT_NAME="${PRODUCT_NAME:-operation-$branch}"
fi

for((i=1;i<=20;i++)); do   
echo $i;
sleep 1
done 
echo "success."

ports="$( docker port operation_parcel-release_release-2.0 )"
success "容器名:        $PRODUCT_NAME"
success "部署路径:      $data_dir"
success "分支:          $branch"
success "端口:"
echo $ports
success "docker-compose.yml: $data_dir/docker-compose.yml"