<h1 class="libTop">xlattice_go</h1>

An implementation of [XLattice](http://xlattice.sourceforge.net)
for the [Go programming language](http://golang.org).

XLattice is a communications library  and a set of tools
for peer-to-peer networks.  The library was originally developed in Java;
more extensive (although somewhat
dated) information on XLattice for Java is available at the
[Java XLattice website](http://www.xlattice.org).

XLattice consists of a number of components.  Generally speaking, those
listed later depend upon some or all of the earlier components.

+ [util](#util)
+ [rnglib](#rnglib)
+ [u](#u)
+ [crypto](#crypto)
+ [transport](#transport)
+ [protocol](#protocol)
+ [overlay](#overlay)
+ [node](#node)
+ [reg](#reg)
+ [httpd](#httpd)

All of these are currently in development.

## <a name="util"></a>util

## <a name="rnglib"></a>rnglib

This version of xlattice_go includes a Go implementation of *rnglib*
for use in testing.  rnglib is a [Go random number generator](rnglib.html)
a drop-in replacement for Go's random number generator.  It

+ is somewhat faster; about 30% in our tests
+ has a number of additional functions for generating random file names,
    directories of random data, etc

## <a name="u"></a>u

and an implementation of **u**, a system for
storing files by their content keys.

[A store organized by content key](u.html)

## <a name="crypto"></a>crypto

## <a name="transport"></a>transport

## <a name="protocol"></a>protocol

### Chunks

[chunks](chunks.html)

## <a name="overlay"></a>overlay

## <a name="node"></a>node

## <a name="reg"></a>reg

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

The xlReg server, its clients, and the cluster members, are all
XLattice [nodes](node.html).

## <a name="httpd"></a>httpd

The go version of XLattice httpd is very much in development.
