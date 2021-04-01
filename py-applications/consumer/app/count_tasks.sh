task_list=($(tkn taskrun list | awk -F' ' '{print $7}'| grep Succeeded))

echo "${#task_list[@]}"