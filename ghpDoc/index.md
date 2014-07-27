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

Most of the XLattice components have been split off as separate Github
repositories.  The project URLs are as shown in the table below.

*component*              | *project*
-------------------------|-------------------------------------------
[util](#util)            | <https://github.com/jddixon/xlUtil_go>
[rnglib](#rnglib)        | <https://github.com/jddixon/rnglib_go>
[u](#u)                  | <https://github.com/jddixon/xlU_go>
[crypto](#crypto)        | <https://github.com/jddixon/xlCrypto_go>
[transport](#transport)  | <https://github.com/jddixon/xlTransport_go>
[protocol](#protocol)    | <https://github.com/jddixon/xlProtocol_go>
[overlay](#overlay)      | <https://github.com/jddixon/xlOverlay_go>
[node](#node)            | <https://github.com/jddixon/xlNode_go>
[nodeID](#nodeID)        | <https://github.com/jddixon/xlNodeID_go>
[reg](#reg)              | <https://github.com/jddixon/xlReg_go>
[httpd](#httpd)          | <https://github.com/jddixon/xlHttpd_go>

All of these are currently in development.

## <a name="util"></a>util

## <a name="rnglib"></a>rnglib

A Go implementation of XLattice's *rnglib*
for use in testing.  rnglib is a [Go random number generator](rnglib.html)
a drop-in replacement for Go's random number generator.  It

+ is somewhat faster; about 30% in our tests
+ has a number of additional functions for generating random file names,
    directories of random data, and such.

## <a name="u"></a>u

This is an implementation of **u**, a system for
storing files by their content keys.

## <a name="crypto"></a>crypto

## <a name="transport"></a>transport

## <a name="protocol"></a>protocol

### Chunks

[chunks](chunks.html)

## <a name="overlay"></a>overlay

## <a name="node"></a>node

## <a name="reg"></a>reg

[xlReg](xlReg.html) is a tool, primarily intended for use in testing,
which facilitates the formation of clusters, groups of cooperating nodes.

On registration, a
client/member is issued a globally unique NodeID, a 256-bit random value.
Once it has an ID, the member can create and/or join clusters.  The cluster has
a maximum size which is set when it is created.  

When members join the cluster they
register their two RSA public keys and either one or two IP addresses.
If the cluster only supports communications between members, members
register only one IP address.  If non-members, *clients*, are allowed to
communicate with the cluster, members register a second address for
that purpose.  When a member has completed registration, it can retrieve
the configuration data other members have registered.

The xlReg server, its clients, and the cluster members, are all
XLattice [nodes](node.html).

## <a name="httpd"></a>httpd

The go version of XLattice httpd is very much in development.
