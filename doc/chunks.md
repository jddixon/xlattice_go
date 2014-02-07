# Chunks

Chunks is a message protocol designed to be intermixed with Protobuf 
[(Google's Protocol Buffer)](http://code.google.com/p/protobuf/) messages.
The latter begin with a varint, a variable-length integer,
which represents the length of the message.
As Protobuf messages have a non-zero length, no message beginning with a 
zero byte can be a valid Protobuf message, so chunks always begin with a 
zero byte, the **magic** field.

In the initial version of Chunks, ***header fields*** are laid out as follows:

|  field | value | bits | bytes |
|:------:|:-----:|:----:|:-----:|
|  magic |   0   |   8  |   1   |
|   type |   0   |   8  |   1   |
|reserved|   0   |  48  |   6   |
| length |   N   |  32* |   4   |
|  index |  ndx  |  32  |   4   |
|  datum | sha3  | 256  |  32   |

The 'datum' field contains the SHA3-256 content hash of the file of 
which this is a chunk.  This is **not** the content hash of the chunk.

The header is followed by `ceil(len/16)` bytes of data, with zero byte
padding added to bring the data length up to a multiple of 16.  That is, 
the data is padded with zeroes as required to bring its length to the next 
whole multiple of 16.  If the length of the data in bytes is already a 
whole multiple of 16, no padding is added.

The chunk data is followed by the 32-byte SHA3-256 hash of the chunk itself, 
where the hash is calculated over the fields named (magic, type, reserved,
length, datum, and data, in that order.

A chunk is assumed to be part of a message.  Chunks are concatenated 
at the destination
in index order to make up the message.  The assumption is that chunks
will normally be transmitted in index order, but that is not necessarily
the case, and using software must tolerate missing, duplicated, and 
out-of-order chunks.

The type 0 Chunk *length* field is constrained to 17 bits and stores 
`N-1`, the length of the chunk less one.   Therefore the 
maximum number of bytes in a chunk is `2^17 = 131072`.  
The upper 15 bits of the length field are reserved and 
must be zero.
