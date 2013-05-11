package xlattice_go

type Key interface {
    Algorithm() string
    GetPublicKey() PublicKey
    GetSigner ()   DigSigner
}
