#!/bin/bash

# for d in *; do
# 	if [ -d $d ]; then
# 	fi
# done

find . -type f -iname 'go.mod' -print0 | xargs -0I% echo "pushd \$(dirname %)>/dev/null && pwd && go mod tidy && popd >/dev/null; echo;echo;echo" | sh
