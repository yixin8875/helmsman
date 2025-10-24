#!/bin/bash

serverName="helmsman"

binaryFile="cmd/${serverName}/${serverName}"

configFile=$1

osType=$(uname -s)
if [ "${osType%%_*}"x = "MINGW64"x ];then
    binaryFile="${binaryFile}.exe"
fi

if [ -f "${binaryFile}" ] ;then
     rm "${binaryFile}"
fi

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

sleep 0.2

go build -o ${binaryFile} cmd/${serverName}/main.go
checkResult $?

trap 'exit 0' SIGINT

# running server
if test -f "$configFile"; then
    ./${binaryFile} -c $configFile
else
    ./${binaryFile}
fi
