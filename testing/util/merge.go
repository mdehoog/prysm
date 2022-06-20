package util

import (
	v2 "github.com/prysmaticlabs/prysm/proto/eth/v2"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// NewBeaconBlockBellatrix creates a beacon block with minimum marshalable fields.
func NewBeaconBlockBellatrix() *ethpb.SignedBeaconBlockBellatrix {
	return HydrateSignedBeaconBlockBellatrix(&ethpb.SignedBeaconBlockBellatrix{})
}

// NewBlindedBeaconBlockBellatrix creates a blinded beacon block with minimum marshalable fields.
func NewBlindedBeaconBlockBellatrix() *ethpb.SignedBlindedBeaconBlockBellatrix {
	return HydrateSignedBlindedBeaconBlockBellatrix(&ethpb.SignedBlindedBeaconBlockBellatrix{})
}

// NewBlindedBeaconBlockBellatrixV2 creates a blinded beacon block with minimum marshalable fields.
func NewBlindedBeaconBlockBellatrixV2() *v2.SignedBlindedBeaconBlockBellatrix {
	return HydrateV2SignedBlindedBeaconBlockBellatrix(&v2.SignedBlindedBeaconBlockBellatrix{})
}

// NewBeaconBlockEip4844 creates a beacon block with minimum marshalable fields.
func NewBeaconBlockEip4844() *ethpb.SignedBeaconBlockWithBlobKZGs {
	return HydrateSignedBeaconBlockEip4844(&ethpb.SignedBeaconBlockWithBlobKZGs{})
}

// NewBlindedBeaconBlockEip4844 creates a blinded beacon block with minimum marshalable fields.
func NewBlindedBeaconBlockEip4844() *ethpb.SignedBlindedBeaconBlockWithBlobKZGs {
	return HydrateSignedBlindedBeaconBlockEip4844(&ethpb.SignedBlindedBeaconBlockWithBlobKZGs{})
}

// NewBlindedBeaconBlockEip4844V2 creates a blinded beacon block with minimum marshalable fields.
func NewBlindedBeaconBlockEip4844V2() *v2.SignedBlindedBeaconBlockEip4844 {
	return HydrateV2SignedBlindedBeaconBlockEip4844(&v2.SignedBlindedBeaconBlockEip4844{})
}
