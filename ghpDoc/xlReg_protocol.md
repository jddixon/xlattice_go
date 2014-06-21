<h1 class="libTop">The XLReg Protocol</h1>

## Registry Credentials

Any XLReg server will profide credentials upon request.
These are conventionally delivered as an ASCII file, `regCred.dat`,
containing

* the registry's name
* its ID as a string of hex digits
* its comms RSA public key in ssh format
* its sig RSA public key in ssh format
* its IP address and port
* and the xlReg protocol version it is running.

An example follows.

<pre><code>regCred {
    Name: xlReg
    ID: 21b95b5c697977d266a0c364e12787ad72bf6fc9346ec0edef351cfb6da90c24
    CommsPubKey: ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCZz5Kdt5XqLRVipmnEEu1eedHmSswP8ZDkbadkRdCrgpGm1OLe79WTrkB0HLW98pjyBooaWLU/thSoB1/2UfkaYdoDHtfHzMKBLUmfR8MCgQaKA3KoOr83wYdtLPYiUmIlg77CjUAuKOPYtd8oy+9TrbM7AwYUZf7Ps/2Lalv7JPQKHX5jyBAjs8nF9LZj+6EhYX0m6RrwyptHjTle7ajQ+6taX+9pZUIY20zu9aiR7j4LNlk2JITOPDk0mr+UsVlI6SfHpuAdy6nsG592bQLT5RF/mD5knh3/EP+b+5yXJHth8myN4UDPIIupinVQ+Vcr0H4y106bebLITWhuJiuN

    SigPubKey: ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDcpMhDQWgLHLcaU4Ed8fHBwLmNOa5RKECmci5VeDczF01R/VaxUcLnna58NM6m1fajNJlS3Z7xICiCwmYFOfJjQ8weuvXebqKUKTZMBghVRJqjPiWGmz9C07U/sTtRrEg0kEUZKepZ6Z9M7VN7eUJwoi+Avp99enTAKmgotYFXn47vpoLDGeKaviHAVcqOHXoQRLfT1Q6vjs/b+yg9lnxRon9kyf3tLopz64Sor6itkI0WhwdWZ0PJHDFW5SfkBhStBW1gC8vED0HO5bbi5iU1NRPiG+nUHm4UYjiQD2DY2PQGXeogZeaqL7ADy8+V0A7TYOkWZTSulK/IuYBY8Clz

    EndPoints {
         TcpEndPoint: 50.18.104.7:55555
    }
    Version: 0.4.3
}
</code></pre>

This particular registry is on stockton.dixons.org, a machine in 
Amazon's AWS/EC2 cloud.  It listens on `55555`, the port conventionally 
used by the xlReg server.

The `Go` version of the xlReg client provides functions to read and write
serialized RegCred files (`xlattice_go.reg.ParseRegCred()` and
`xlattice_go.reg.String()` respectively.

## Hello and Reply

All communications with the XLReg server must begin with a Hello/Reply
exchange, which verifies to the client the identity of the server and
establishes the AES IV and key used to encrypt further message
exchanges between the XLReg server and client.

### Version Numbers

XLReg version numbers are 4-byte little-endian values.  These may be
thought of as `a.b.c.d`, where `a`, `b`, `c`, and `d` are unsigned byte values.
`d` is the lowest value and so will appear first in wire format.  When
serialized as strings, version numbers are conventionally written in
big-endian form, except that if the low-order fields are all zero, only
the higher-order fields appear, so `1.2.0.0` will normally be written `1.2`
(but will appear on the wire as `0.0.2.1`.  Similarly `1.2.3.0` serializes
as `1.2.3` but appears on the wire as `0.3.2.1` and `1.2.3.4` appears on
the wire as `4.3.2.1.

### Hello

XLReg communications begin when one machine, the **client**,
sends a Hello message to another, the **server** in this context.  The
Hello message consists of

* a 16-byte AES initialization vector (IV)
* a 32-byte AES key
* an 8 byte salt (random value)
* a 4-byte little-endian version number

AES is a block cipher used for high speed encryption.  AES has a standard
block size of 16 bytes.  The key is a whole number of such blocks.  The IV is
used to set up cipher-block chaining (CBC), a mode of operation in which
each block of plaintext is XORed with the previous block of ciphertext
before encryption.  The first block of plaintext has no previous block
of ciphertext, so an IV is used instead.

Clients should encrypt the message using RSA-OAEP, SHA1, and a 20-byte
random value, the oaepSalt.  No label is used.  In Go this is done using

	ciphertext,err = rsa.EncryptOAEP(sha, oaepSalt, ck, data, nil)

where `sha` represents an instance of the SHA1 hash function and `ck` is
a pointer to an RSA public key.

In production all of the random values (IV, key, salt, and oaepSalt)
should be generated using a
secure, crypto-quality random number generator.

The Hello message is encrypted using **the server's** public key.
The message to be encoded must not be longer than the length of modulus
less twice the hash length plus 2.  As we use SHA1 the hash length is
20 bytes, so for a 1024-bit = 128 byte RSA key the maximum message length
is

	128 - (2 * 20 + 2) = 86 bytes

The hello message is 60 bytes long and so fits.  We recommend using 1024-bit
RSA keys for testing but 2048-bit or larger in production.

The Hello message is sent to the server over a TCP connection using its
well-known address.  Conventionally the xlReg server listens on port 55555.

### Reply

The server decrypts the message using its RSA private key - and the
server's correctly encrypted reply proves to
the client that the server has that private key.  The server examines
the message on receipt;
if it is not well-formed the server silently discards it.  Otherwise the
server encrypts its reply using the AES IV and key from the Hello message.
The server's reply consists of

* iv2, the 16-byte IV to be used in further communications
* key2, the 32-byte AES key to be used in such communications
* salt2, an 8-byte random value
* salt1, the salt from the Hello message
* version2, a 4-byte little endian value

The first three fields should be generated by a secure random number generator.

The version number in the Hello message is a _proposed_ version number,
the version preferred by the client.  The server will reply with the
version number that it prefers, and this is the protocol version that
will be used in subsequent messages.  In this implementation the server
simply ignores the version proposed by the client.

If the reply from the server cannot be decoded, or if salt1 does not
match the value in the Hello message, the client should silently close
the connection.

Otherwise iv2 and key2 will be used in further messages between this
client and server.

## The Role of Protocol Buffers

All xlReg client-server sessions must begin with a Hello/Reply sequence,
which establishes the identify of the server and determines the AES IV
and key used to encrypt all further messages in the session.  While the
Hello/Reply sequence is specified in terms of a pattern of bits on the
wire, ClientMsg and OKMsg are specified by a
[Google Protocol Buffers](http://code.google.com/p/protobuf)
protocol description file, `p.proto`.  This is used to generate libraries
specific to the particular language.

Any particular XLRegMsg message is first translated into wire format
by a Protobuf library call, then PKCS7-padded to a whole number of
16-byte AES blocks, and then AES-encrypted using the IV and key set
during the Hello/Reply sequence.  In Go this is done by a call to

    EncodePadEncrypt(msg *XLRegMsg, engine cipeher.BlockMode)

which returns either a byte slice or an error.

The receiver inverts this process.  It gets a byte slice off the wire
and makes a call to

    DecryptUnpadDecode(ciphertext []byte, engine cipher.BlockMode)

which returns either a pointer to an XLRegMsg or an error.

## XLReg Protobuf Protocol Description

[xlReg Protobuf protocol](xlReg_protobuf.html)

## Client and OK

The ClientMsg, like all other XLRegMsg types, can only be sent to
the server after AES encryption is set up by the Hello/Reply sequnce.
The message descriptions that follow are expressed in terms of the
Protobuf message spec.

### Client Message

The formal spec provides two versions of the Client message, one
containing a Token and the other only the ClientID.  At this time
only the token-based message should be used.

The token embedded in the client message consists of

* the client **Name**, whose leading character should be a letter; other
  characters should be alphanumeric
* an unsigned 64-bit **Attrs** field; this is the client's proposed
  value (and is ignored in the current implementation)
* a serialized RSA public key, the **CommsKey**, used for encrypting
  (small) messages
* another serialzied RSA public key, the **SigKey**, used only for
  creating digital signatures
* an array of strings, **MyEnds**, which must be the serialized TCP/IP
  address, including port number, that the prospective member listens on
* and finally **DigSig**, the RSA/SHA1 digital signature over each of
  the preceding fields, in order

Note that the ClientID should **not** be included in the token; it is
assigned by the server in response to this message.

### OK Message

The server examines the client message on receipt.  If the message is
ill-formed, the server simply closes the connection and discards the message.
Otherwise it determines the value of the **Attrs** field for the client
(currently it just accepts the client's proposed value and returns it) 
and constructs a unique random ID for the client.  This is a 
256-bit / 32-byte **ClientID**.  The AES-encrypted reply to the client 
contains both of these fields.

* **ClientID**, a 32-byte byte array
* **Attrs**, an unsigned 64-bit int

## Create and CreateReply

The Create/CreateReply sequence is necessary only if the cluster does not
already exist.

The client creating the cluster need not be one of the members.  It is
in fact convenient to use a dedicated admin client to create the cluster.

### Create Message

A cluster definition must be sent to the registry for each cluster
once and only once.  A second cluster definition with the name of an
existing cluster will be rejected by the server, which will reply
with an XLReg Error message rather than a CreateReply.

The cluster Create message consists of

*  **ClusterName**, which should be an acceptable xlReg name, and so
   begin with a letter (`[a-zA-Z_]`) and otherwise alphanumeric
   (`[a-zA-Z_0-9\`).
* **ClusterAttrs**, an unsigned 64-bit integer, which currently is ignored
* **ClusterSize**, the maximum number of members in the cluster; this is
  a 32-bit integer constrained to be in the range 1..64 inclusive
* **EndPointCount**, which in the current implementation should be 1 or 2

The cluster name must not already be in use. If it is, the server will 
reply with an error message and close the connection.

The `EndPointCount` is the number of endpoints that each cluster member
must have.  A member must listen on each address/port number listed.
Conventionally `endPoint[0]` is used for intra-cluster communications,
communications between cluster members, and `endPoint[1]`, if defined,
is used by clients to reach cluster members.

### CreateReply Message

The `CreateReply` message is sent by the server to the client to confirm
cluster attributes.  It consists of

* **ClusterID**, wich is assigned by the server
* ** ClusterSize***, the maximum number of members, which may have been
  adjusted by the server
* **EndPointCount***, which also is subject to change by the server

### Error Message

In the current implementation this message has a single field, an
`ErrDesc` string.  If the server sends the client an error message,
it then closes the connection.

## Join and JoinReply

### Join Message

In the current implementation a Join message has only a single field,
`ClusterName`, a string.  If the cluster exists and few members have
joined than the cluster size, the xlReg server will sdnd a JoinReply
message.  Otherwise, the server will send an error message and close
the connection.

### JoinReply Message

If a Join request succeeds, the server responds with a JoinReply
containing:

* the **EndPointCount**, a possibly adjusted endPoint count and
* **ClusterID**, the usual 256-bit binary value

This is the point at which prospective members knowing only the
cluster name will have learned the cluster ID.

## GetCluster and ClusterMembers

A client sends the xlReg server GetCluster messages until it has
information on all members or until some limit (MAX_GET in the Go code)
has been exceeded.  The client does not have to be a member in order
to make this request.

### GetCluster Message

The message has two fields:

* **ClusterID**, the 32-byte cluster ID, and
* **Which**, a bit map identifying which members information is needed about

The bit map is a little-endian sequence of 64 bits, with the low-order
bits in each byte mapping to 0, 1, 2, and so forth.  The server ignores
any bits that are out of range, so `0xffffffff` simply means "all members".
The client maintains a local bit map identifying those members that it
has information on.  It will loop until it has information on all members
or until MAX_GET is exceeded.

### ClusterMembers Message

The xlReg server replies to the GetCluster message with a ClusterMembers
message consisting of

* **ClusterID**, the usual 256-bit cluster identifier
* **Which**, a bit map identifying which members the message contains
    information about, and
* **Tokens**, an array of tokens as described above under the Client
    Message heading

The tokens are in the same order as bits in the bit map, so the client
iterates through the bit map and copies each token to the corresponding
slot in its own members table.

If a cluster member is making this request -- that is, if the client has
already joined the cluster -- the reply will normally contain a token
for the requesting client.

The xlReg server does not maintain state regarding which membership
information the client has collected.  If asked repeatedly for the same
information, it will send the information to the client repeatedly.

## Bye and Ack

### Bye Message

The **Bye** is sent by the client to signal that the session is over.
It contains no further information.

### Ack Message

The server's reply is a simple **Ack**.  It contains no further
information.  After sending it, the server closes the connection.  On
receiving an the Ack, the client does the same.

