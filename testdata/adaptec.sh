#!/usr/bin/env bash

if [[ $1 = "list" ]]; then cat testdata/adaptec/controllers.txt
elif [[ $1 = "getconfig" ]] && [[ $3 = "ld" ]]; then cat testdata/adaptec/logicaldrives.txt
elif [[ $1 = "getconfig" ]] && [[ $3 = "pd" ]]; then cat testdata/adaptec/physicaldrives.txt
elif [[ $1 = "getconfig" ]] && [[ $3 = "ad" ]]; then cat testdata/adaptec/controllerStatus.txt
fi
