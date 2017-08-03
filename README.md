Fence Executor
==============
(This project is based on a fork of sgotti/go-fence project)

Fence Executor is a golang fencing server pluggable with different fence providers.
Fencing is the ability to perform power management actions to isolate a node.

## Fence Providers

At the moment there's one fencing provider that uses the fence agents provided by the redhat fence agents project (used by redhat cluster/pacemaker, oVirt and other projects).

## Api

POST https://localhost:7777/fence

'''
{
    "script": "fence_apc_snmp"
    "address": "fqdn.com"
    "password": "password"
    "username": "user"
    "secure": False
    "options": {"port": 21,}
}
'''
