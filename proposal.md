

# Proposal: Managing Borah Specifications and Attribution
Author(s): Kyle Shannon <kyle@pobox.com>

Last updated: 2020-10-22

## Abstract

Research Computing Services (RCS) asked researchers to give proper attribution
when HPC resources are used for research.  There are various methods to achieve
this, including issuing a DOI for the resource, or specifying text for authors
to use in the publication.

This proposal defines a standard way to handle the attribution, and the
underlying specifications that may change over time.

## Background

The previous instantiation of Boise State HPC resources (R2) has a DOI on the
library's publication archive (10.18122/B2S41H), and the appropriate text for
attribution is defined [here](https://www.boisestate.edu/rcs/publications-acknowledgment/).

Issues arose when the specification of R2 changed.  This usually happens when
resources are added, such as condominium nodes, or infrastructure such as login
nodes.  The changes had to be reflected in the DOI indexed document on
scholarworks, which in turn requires a new DOI.

The changing of the DOI is undesirable, and can cause ambiguity.  An attempt to
handle this was implemented in github and zenodo, which tracked changes in the
specifications, but still required several DOIs.

Most users do not need to reference the exact specification of the compute
cluster, and those who do can be provided with a secondary service to lookup
specifications for a given date.  This isn't an imaginary issue, as we have
researchers on campus that manipulate machine code via compiler components,
which can heavily depends on specific architectures.

## Proposal

We propose using two systems to track the referencing and the specification
tracking for the compute cluster.  Between the two, we should be able to
pinpoint near-exact specifications for a given date, while not changing the DOI
to be referenced in scholarworks.

The specifications over time will be stored in RCS github repository, and
updated in a timely manner as resources are added or removed.  This will give
us an effective timeline of what the cluster looks like at any point in time.

A general description of the cluster will be submitted to the library for
scholarworks, and it will reference the github repository, as well as the
specification page on the RCS website.  This is an attempt at providing a
solution for both the single DOI and the fine-grained specs when needed.

## Rationale

*TODO(KYLE): GO OVER WITH MENDI...*

RCS needs a _simple_ way to allow for attribution, while keeping a valid
history of the cluster's specifications.  

### DOI for every spec change

This method is difficult to assign attribution, as authors may not always know
which DOI to use.

## Implementation

The github repository will contain a group of [yaml](https://yaml.org/) files
that represent all nodes and some general hardware on the cluster.  `borah.yml`
will hold the cluster-wide and general information such as interconnect and
storage technology.  All other yaml files contain a list of node definitions
that represents a homogeneous set of nodes.  There can be several node
defintions in file.  The files are consumed as if all yaml files are appended
to the `borah.yml` file.  Separate files are a matter of convenience.

The `borah.yml` file has the general structure:

```
  # Meta-information for the entire Borah cluster
  name: Borah
  interconnect: Mellanox Non-Blocking HDR200/HDR100 Infiniband
  internet: 100GB ethernet
  # DO NOT EDIT BELOW THIS LINE
  nodes:
```

The `nodes:` tag _must not_ be removed, as all other yaml files are structured
to fit under that tag.  Any cluster-wide specification can go in the
`borah.yml` file, above the `nodes:` tag.

The node specification files consist of lists of nodes.  A node definition
consists of:

```
  - type: Compute
    make: Dell
    model: C6420
    cpu: Intel Xeon Gold 6252
    cpus: 2
    cpucores: 24
    ram: 192
    count: 40
    owner: Boise State
  ```

another example, perhaps stored in `condo.yml`:

```
  - type: Condo
    make: Dell
    model: C6420
    cpu: Intel Xeon Gold 6252
    cpus: 2
    cpucores: 24
    ram: 192
    count: 1
    owner: Eric Hayden
  - type: Condo
    make: Dell
    model: C6420
    cpu: Intel Xeon Gold NEW
    cpus: 2
    cpucores: 20
    count: 2
    owner: kyle
```

These files can be ingested to generate a variety of text, tables, etc for
display on websites or in reports, as well as providing a history of the
cluster.

Preliminary data files can be found [here](https://github.com/bsurc/borah-specs).
