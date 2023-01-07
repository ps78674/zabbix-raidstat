#!/usr/bin/env bash

if [[ $1 = "-AdpGetPciInfo" ]]; then cat testdata/megacli/controllers.txt
elif [[ $1 = "-LdInfo" ]] && [[ $3 = "-Lall" ]]; then cat testdata/megacli/logicaldrives.txt
elif [[ $1 = "-PDList" ]]; then cat testdata/megacli/physicaldrives.txt
elif [[ $1 = "-AdpAllInfo" ]]; then cat testdata/megacli/controllerStatus.txt
elif [[ $1 = "-AdpBbuCmd" ]] && [[ $3 = "-GetBbuStatus" ]]; then cat testdata/megacli/controllerBBUStatus.txt
elif [[ $1 = "-LdInfo" ]] && [[ $3 != "-Lall" ]]; then cat testdata/megacli/logicaldriveStatus.txt
elif [[ $1 = "-pdInfo" ]]; then cat testdata/megacli/physicaldriveStatus.txt
fi
