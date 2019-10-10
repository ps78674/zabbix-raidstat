### Zabbix RAID monitoring for Adaptec and HP Smart Array
Simple parser for `arcconf` and `ssacli` written in Go.

Zabbix template provides discovery for controllers, logical and physical drives.
![Discovery](https://user-images.githubusercontent.com/31385755/65332764-f9f3f380-dbc7-11e9-9d08-9a2e5bc236bf.png)

Configured host must have macros {$RAID_VENDOR} as values for cli option `-vendor`.
![Example host](https://user-images.githubusercontent.com/31385755/65949183-5cf54e00-e444-11e9-9070-ef570a53c7e4.png)

```
Usage of ./raidstat:
  -d string
     Discovery option, one of 'ct | ld | pd'
  -indent int
     Indent JSON output for 
  -path string 
     RAID tool full path, like '/opt/arcconf'
  -s string
     Status option, one of 'ct,<CONTROLLER_ID> | ld,<CONTROLLER_ID>,<LD_ID> | pd,<CONTROLLER_ID>,<PD_ID>'
  -vendor string
     RAID tool vendor, one of 'adaptec | hp'
```
If flag '-path' is not set, then tool uses filename from map.
```
vendorTools = map[string]string{
    "adaptec": "arcconf",
    "hp":      "ssacli",
}

...

if len(toolBinary) == 0 {
    toolBinary = vendorTools[toolVendor]
}
```
Instead, userparameters can be modified lke this:  
UserParameter=raidstat.discovery.controllers[*], sudo /opt/raidstat/raidstat -vendor $1 -path <PATH_TO_RAID_TOOL> -d ct

## Installation:

1. Provide `zabbix_agentd` process user passwordless sudo access to raidstat binary - see `raidstat/zabbix/raidstat.sudoers`
2. Copy `zabbix/userparameter_raidstat.conf` to `/etc/zabbix/zabbix_agentd.d`
3. Copy compiled binaries to `/opt/raidstat`
4. Import template`zabbix/zbx_raid_monitoring.xml`
