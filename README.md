Fence Executor
==============
(This project is based on a fork of sgotti/go-fence project)

## General
Fencing is the ability to perform power management actions to isolate a node in cluster.

# Usage
Fence Executor is a golang fencing executable that can be used to perform fence operation. The implementation is pluggable and supports providers and agents
'''
go run main.go [fqdn] [username] [password] [plug(on\off)] [provider(redhat)] [agent(fence_apc_snmp)]
'''
