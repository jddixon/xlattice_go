REGISTRY
========

As the term is used here, *registry* is a protocol for communicating
between parties in a traditional server/client relationship.  On the
one hand we have independent nodes which may cooperate as a cluster,
a collection of peers that communicate with one another over links
dedicated to this purpose.  On the other hand we have a single server
which is identified by its nodeID (a 20- or 32-byte value).  The 
server has a well-known address, a tcp/ip port in an address region,
an overlay, shared by the nodes wishing to form a cluster.  It also
has a well-known RSA public key.

Clients join a cluster by sending a Hello message to the server to agree
on communications and then send an encrypted Join message with their 
details.  

The Hello is encrypted using the server's RSA public key.  It contains a 
salt ("salt1") and an AES IV and key ("KeyIV") used only to encrypt the 
reply, the HelloReply.

The HelloReply contains the original salt, another salt ("salt2"), and
the AES IV and key (KeyIV2) used for the rest of the session.   By 
deciphering the first salt, the server has proved its identity - that is,
it has proved that it has the secret RSA key corresponding to the public
key.  If the Hello is in some way ill-formed, the server will silently 
close the connection.  The client will do the same if the HelloReply is 
not properly encrypted using KeyIV or does not contain the correct value 
for the salt.  From this point KeyIV2 is used by both sides to encrypt 
session traffic.

