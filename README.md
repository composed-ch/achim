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
            - [ ] handle cloud-init data
        - ~~[ ] check (check-state)~~ (handled by list)
        - [x] deprotect
        - [ ] destroy
        - [x] label (label-all-instances)
        - [x] list (list-instances)
        - [ ] probe
        - [x] protect
        - [ ] resize (resize-disk)
        - [ ] scale (scale-instance)
        - [x] start
        - [x] stop
        - [x] type (list-instance-types)
    - group
        - [ ] create (create-group)
        - [ ] export-overview (export-group-overview)
        - [ ] export-inventory (export-inventory)
        - [ ] export-playbook (export-user-playbook)
        - [ ] file-from-text (TBD)
    - network
        - [ ] attach
        - [ ] create (create-network)
        - [ ] cleanup (cleanup-network)
        - [ ] flush (flush-networks)
        - [ ] list (list-network)
        - [ ] destroy (destroy-network)
    - scenario
        - [ ] create (create-scenario)
        - [ ] destroy (destroy-scenario)
        - [ ] export-overview (export-scenario-overview)
    - dns
        - [ ] flush
        - [ ] sync (sync-dns)
    - image
        - [x] list (list-images)
- other tasks
    - [ ] replace `%v` with `%w` for proper error wrapping