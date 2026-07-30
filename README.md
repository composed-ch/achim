# achim

Advanced Cloud Hyperscaling Infrastructure Manager (reimplementation in Go) using the [urfave/cli](https://cli.urfave.org/) and [Egoscale](https://github.com/exoscale/egoscale) libraries.

## Status

This Go implementation is a port from the [Python implementation](https://github.com/patrickbucher/achim)

- command structure (and existing achim command, if different)
    - tenant
        - [x] add ([new])
        - [x] default ([new])
        - [x] remove ([new])
    - instance
        - [x] create (create-instance)
        - [x] deprotect
        - [x] destroy
        - [x] label (label-all-instances)
        - [x] list (list-instances)
        - [x] overview (export-group-overview)
        - [x] probe
        - [x] protect
        - [x] ~~resize~~ embiggen (resize-disk)
        - [x] scale (scale-instance)
        - [x] start
        - [x] stop
        - [x] type (list-instance-types)
    - group (high priority)
        - [x] create (create-group)
            - [ ] handle cloud-init data
        - [ ] export-inventory (export-inventory)
            - TODO: part of instance or group?
        - [ ] export-playbook (export-user-playbook)
            - TODO: part of instance or group?
        - [x] file-from-text
    - dns (high priority)
        - [ ] flush
        - [ ] sync (sync-dns)
    - network (lower priority)
        - [ ] attach
        - [ ] create (create-network)
        - [ ] cleanup (cleanup-network)
        - [ ] flush (flush-networks)
        - [ ] list (list-network)
        - [ ] destroy (destroy-network)
    - scenario (lower priority)
        - [ ] create (create-scenario)
        - [ ] destroy (destroy-scenario)
        - [ ] export-overview (export-scenario-overview)
    - image
        - [x] list (list-images)
- other tasks
    - [ ] replace `%v` with `%w` everywhere for proper error wrapping
    - [ ] consider `--verbose`/`-v` (or `--silent/-s`?) flag for certain commands
        - alternative: `--dry`/`-d` run flag (what _would_ happen _if_?)
        - instance create/destroy, start/stop, protect/deprotect, resize/scale
        - consider making `--by`/`-b` argument optional for some commands (e.g. `instance list` and the like)