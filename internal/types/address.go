package types

type AddrRangeIterable interface {
	GetNext() *PublicDerivationAddressRangeData
	GetRangesSize() uint32
}

type PublicDerivationAddressRangeData struct {
	AccountIndex     uint32
	InternalIndex    uint32
	AddressIndexFrom uint32
	AddressIndexTo   uint32
	AddressIndexDiff int32
}
