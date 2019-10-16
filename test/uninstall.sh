#!/bin/bash  
# 测试用的脚本
success() { printf "${green}✔ %s${reset}\n" "$@"
}

for((i=1;i<=20;i++));  
do   
echo $i;  
sleep 1
done 
echo "success."

success "success."