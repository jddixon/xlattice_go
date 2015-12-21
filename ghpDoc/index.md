<h1 class="libTop">xlattice_go</h1>

An implementation of the open source
[XLattice](http://jddixon.github.io/xlattice)
project for the
[Go programming language](http://golang.org).

## Project Description

XLattice is a communications library and a set of tools
for peer-to-peer networks.  The library was originally developed in Java;
more extensive (although somewhat
dated) information on XLattice for Java is available at the
[Java XLattice website](http://www.xlattice.org).

## Components

XLattice consists of a number of components.  Generally speaking, those
listed later depend upon some or all of the earlier components.

| component               | Go project documentation                   |
|-------------------------|--------------------------------------------|
| [rnglib](#rnglib)       | <https://jddixon.github.io/rnglib_go>      |
| [util](#util)           | <https://jddixon.github.io/xlUtil_go>      |
| [u](#u)                 | <https://jddixon.github.io/xlU_go>         |
| [crypto](#crypto)       | <https://jddixon.github.io/xlCrypto_go>    |
| [transport](#transport) | <https://jddixon.github.io/xlTransport_go> |
| [protocol](#protocol)   | <https://jddixon.github.io/xlProtocol_go>  |
| [overlay](#overlay)     | <https://jddixon.github.io/xlOverlay_go>   |
|                         | <https://jddixon.github.io/xlNodeID_go>    |
| [node](#node)           | <https://jddixon.github.io/xlNode_go>      |
| [cluster](#cluster)     | <https://jddixon.github.io/xlCluster_go>   |
| [reg](#reg)             | <https://jddixon.github.io/xlReg_go>       |
| [httpd](#httpd)         |                                            |

All of these are currently in development.

### <a name="rnglib"></a>rnglib

This version of xlattice_go includes a Go implementation of
the python package
[rnglib](https://jddixon.github.io/rnglib)
for use in testing. rnglib_go is a [Go random number generator](rnglib.html)
a drop-in replacement for Go's random number generator.  It

* is somewhat faster; about 30% in our tests
* has a number of additional functions for generating random file names,
    directories of random data, etc

### <a name="util"></a>util

### <a name="u"></a>u

and an implementation of **u**, a system for
storing files by their content keys.

[A store organized by content key](u.html)

### <a name="crypto"></a>crypto

### <a name="transport"></a>transport

### <a name="protocol"></a>protocol

#### Chunks

[chunks](chunks.html)

### <a name="overlay"></a>overlay

### <a name="node"></a>node

### <a name="reg"></a>reg

[xlReg](https://jddixon.github.io/xlReg_go)
is a tool, primarily intended for use in testing,
which facilitates the formation of clusters, groups of cooperating nodes.
On registration, a
client/member is issued a globally unique NodeID, a 256-bit random value.
Once it has an ID, the member can create and/or join clusters.  The cluster has
a maximum size set when it is created.  When members join the cluster they
register their two RSA public keys and either one or two IP addresses.
If the cluster only supports communications between members, members
register only one IP address.  If non-members, clients, are allowed to
communicate with the cluster, members register a second address for
that purpose.  When a member has completed registration, it can retrieve
the configuration data other members have registered.

The xlReg server, its clients, and the cluster members, are all XLattice
[nodes](https://jddixon.github.io/xlattice/node).

### <a name="httpd"></a>httpd

## Project Status

The Go version of XLattice httpd is very much in development.

