# xlattice_go

An implementation of [XLattice](http://xlattice.sourceforge.net)
for the Go language.  XLattice is a communications library 
for peer-to-peer networks.  More extensive (although somewhat
dated) information on XLattice is available as the 
[XLattice website](http://www.xlattice.org).

This version of xlattice-go includes a Go implementation of *rnglib*
for use in testing and an implementation of **u**, a system for
storing files by their content keys.

## xlReg

[xlReg](xlReg.html) is a tool, primarily intended for use in testing, 
which facilitates the formation of clusters.  On registration, a 
client/member is issued a globally unique NodeID, a 256-bit random value.  
The member can then create and/or join clusters.  The cluster has
a maximum size set at creation.  When members join the cluster they
register their two RSA public keys and either one or two IP addresses.
If the cluster only supports communications between members, members
register only one IP address.  If non-members, clients, are allowed to
communicate with the cluster, members register a second address for 
that purpose.  When a member has completed registration, it can retrieve
the configuration data other members have registered.

The xlReg server, its clients, and the cluster members, and all 
XLattice [nodes](node.html).

## Other Protocols

### Chunks

[chunks](chunks.html)

## u

[A store organized by content key](u.html)

## rnglib

[A Go random number generator](rnglib.html)
