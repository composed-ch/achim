# achim

Advanced Cloud Hyperscaling Infrastructure Manager (reimplementation in Go) using the [Cobra](https://github.com/spf13/cobra) and [Egoscale](https://github.com/exoscale/egoscale) libraries.

## Status

This Go implementation is a port from the [Python implementation](https://github.com/patrickbucher/achim)

- command structure (and existing achim command, if different)
    - access
        - [ ] add ([new])
        - [ ] remove ([new])
    - instance
        - [ ] create (create-instance)
        - [ ] check (check-state)
        - [ ] protect
        - [ ] deprotect
        - [ ] destroy
        - [ ] label (label-all-instances)
        - [ ] list (list-instances)
        - [ ] probe
        - [ ] resize (resize-disk)
        - [ ] scale (scale-instance)
        - [ ] start
        - [ ] stop
    - group
        - [ ] create (create-group)
        - [ ] export-overview (export-group-overview)
        - [ ] export-inventory (export-inventory)
        - [ ] export-playbook (export-user-playbook)
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
    - images
        - [ ] list (list-images)
        - [ ] list-types (list-instance-types)

