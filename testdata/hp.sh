#!/usr/bin/env bash

if [[ $@ = "ctrl all show" ]]; then cat testdata/hp/controllers.txt
elif [[ $1 = "ctrl" ]] && [[ $3 = "ld" ]] && [[ $4 = "all" ]] && [[ $5 = "show" ]]; then cat testdata/hp/logicaldrives.txt
elif [[ $1 = "ctrl" ]] && [[ $3 = "pd" ]] && [[ $4 = "all" ]] && [[ $5 = "show" ]]; then cat testdata/hp/physicaldrives.txt
elif [[ $1 = "ctrl" ]] && [[ $3 = "show" ]] && [[ $4 = "status" ]]; then cat testdata/hp/controllerStatus.txt
elif [[ $1 = "ctrl" ]] && [[ $3 = "ld" ]] && [[ $5 = "show" ]] && [[ $6 = "detail" ]]; then cat testdata/hp/logicaldriveStatus.txt
elif [[ $1 = "ctrl" ]] && [[ $3 = "pd" ]] && [[ $5 = "show" ]] && [[ $6 = "detail" ]]; then cat testdata/hp/physicaldriveStatus.txt
fi
