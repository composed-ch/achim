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
        - [x] inventory (export-inventory)
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
    - group
        - [x] create (create-group)
        - [x] file-from-text
        - [x] playbook (export-user-playbook)
    - dns
        - [x] flush
        - [x] sync (sync-dns)
    - network
        - [x] attach
        - [x] create (create-network)
        - [x] cleanup (cleanup-network)
        - [x] flush (flush-networks)
        - [x] list (list-network)
        - [x] destroy (destroy-network)
    - scenario
        - [x] create (create-scenario)
            - TODO: network labels (if needed)
            - TODO: cloud init data (if needed)
        - [ ] export-overview (export-scenario-overview)
    - image
        - [x] list (list-images)