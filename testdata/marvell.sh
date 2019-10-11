#!/usr/bin/env bash

if [[ $1 = "info" ]] && [[ $2 = "-o" ]] && [[ $3 = "hba" ]]; then cat testdata/marvell/controllers.txt
elif [[ $1 = "info" ]] && [[ $2 = "-o" ]] && [[ $3 = "ld" ]]; then cat testdata/marvell/logicaldrives.txt
elif [[ $1 = "info" ]] && [[ $2 = "-o" ]] && [[ $3 = "pd" ]]; then cat testdata/marvell/physicaldrives.txt
fi
