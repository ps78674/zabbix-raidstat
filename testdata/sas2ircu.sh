#!/usr/bin/env bash

if [[ $1 = "list" ]]; then cat testdata/sas2ircu/list.txt
elif [[ $1 = "0" ]] && [[ $2 = "status" ]]; then cat testdata/sas2ircu/status.txt
elif [[ $1 = "0" ]] && [[ $2 = "display" ]]; then cat testdata/sas2ircu/display.txt
fi
