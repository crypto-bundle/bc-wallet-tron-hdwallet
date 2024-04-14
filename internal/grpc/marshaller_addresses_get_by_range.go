package grpc

import (
	pbCommon "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/common"
)

type marshallerAddressesGetByRange struct {
	addressesIdentities []*pbCommon.DerivationAddressIdentity
}

func (m *marshallerAddressesGetByRange) MarshallPath(accountIndex, internalIndex, addressIdx,
	position uint32,
	address string,
) {
	addressEntity := &pbCommon.DerivationAddressIdentity{}

	addressEntity.AccountIndex = accountIndex
	addressEntity.InternalIndex = internalIndex
	addressEntity.AddressIndex = addressIdx
	addressEntity.Address = address

	m.addressesIdentities[position] = addressEntity

	return
}

func (m *marshallerAddressesGetByRange) GetMarshaled() []*pbCommon.DerivationAddressIdentity {
	return m.addressesIdentities
}

func newAddrRangeMarshaller(addListSize uint32) *marshallerAddressesGetByRange {
	return &marshallerAddressesGetByRange{
		addressesIdentities: make([]*pbCommon.DerivationAddressIdentity, addListSize),
	}
}
