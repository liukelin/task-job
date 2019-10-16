#/bin/bash
cd `dirname $0`
mode=$1
case "$mode" in
    'windows')
        ;;
    'linux')
        ;;
    'darwin')
        ;;
    'freebsd')
        ;;
    *)
        $mode = "linux"
        ;;
esac

CGO_ENABLED=0 GOOS=$mode GOARCH=amd64 go build ./main.go

if [[ "$1" == "windows" ]];then
    mv ./main.exe ./windows.exe
else
    mv ./main ./$1
fi

echo "$1 build success."

