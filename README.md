### Zabbix RAID monitoring for Adaptec/Microsemi, HP Smart Array and Lenovo M.2 RAID (mvcli)
Simple parser for `arcconf`, `ssacli` and `mvcli` written in Go.

Zabbix template provides LLD for controllers, logical and physical drives.
![Discovery](https://user-images.githubusercontent.com/31385755/65332764-f9f3f380-dbc7-11e9-9d08-9a2e5bc236bf.png)

Configured host must have macros {$RAID_VENDOR} (as value for cli option `-vendor`).
![Example host](https://user-images.githubusercontent.com/31385755/65949183-5cf54e00-e444-11e9-9070-ef570a53c7e4.png)

```
Usage of ./raidstat:
  -d string
     Discovery option, one of 'ct | ld | pd'
  -indent int
     Indent JSON output for 
  -s string
     Status option, one of 'ct,<CONTROLLER_ID> | ld,<CONTROLLER_ID>,<LD_ID> | pd,<CONTROLLER_ID>,<PD_ID>'
  -vendor string
     RAID tool vendor, one of 'adaptec | hp | marvell'
```
Config file `config.json` is used for raid vendors -> tools configuration.
```
{
    "vendors": {
        "hp": "ssacli",
        "vendor1": "/PATH/TO/BINARY1",
        "vendor2": "/PATH/TO/BINARY2"
    }
}
```
Vendor name is used as plugin name (like "hp.so").

## Installation:

1. Provide `zabbix_agentd` process user with passwordless sudo access to raidstat binary - `raidstat/zabbix/raidstat.sudoers`
2. Copy `zabbix/userparameter_raidstat.conf` to `/etc/zabbix/zabbix_agentd.d`
3. Copy compiled binaries to `/opt/raidstat`
4. Import template`zabbix/zbx_raid_monitoring.xml`
