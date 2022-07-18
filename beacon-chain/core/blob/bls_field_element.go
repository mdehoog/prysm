package blob

import (
	"github.com/holiman/uint256"
	v1 "github.com/prysmaticlabs/prysm/proto/engine/v1"
	"math/big"
)

const blsModulusStr = "52435875175126190479447740508185965837690552500527637822603658699938581184513"

var blsModulusBig = big.Int{}
var blsModulusInt = uint256.Int{}

func init() {
	if err := blsModulusBig.UnmarshalText([]byte(blsModulusStr)); err != nil {
		panic(err)
	}
	if overflow := blsModulusInt.SetFromBig(&blsModulusBig); overflow {
		panic("overflow")
	}
}

type BLSFieldElement uint256.Int

func (e *BLSFieldElement) SetFromBigMod(b *big.Int) bool {
	i := &big.Int{}
	i.Mod(b, &blsModulusBig)
	return (*uint256.Int)(e).SetFromBig(i)
}

func (e *BLSFieldElement) SetBytesMod(b []byte) *BLSFieldElement {
	x := (*uint256.Int)(e)
	x.SetBytes(b)
	return (*BLSFieldElement)(x.Mod(x, &blsModulusInt))
}

func (e *BLSFieldElement) MulMod(b *BLSFieldElement) *BLSFieldElement {
	x := (*uint256.Int)(e)
	y := (*uint256.Int)(b)
	return (*BLSFieldElement)(x.MulMod(x, y, &blsModulusInt))
}

func (e *BLSFieldElement) AddMod(b *BLSFieldElement) *BLSFieldElement {
	x := (*uint256.Int)(e)
	y := (*uint256.Int)(b)
	return (*BLSFieldElement)(x.AddMod(x, y, &blsModulusInt))
}

func (e *BLSFieldElement) Bytes() []byte {
	return (*uint256.Int)(e).Bytes()
}

func BlobsToBLSFieldElements(blobs []*v1.Blob) [][]BLSFieldElement {
	a := make([][]BLSFieldElement, len(blobs), len(blobs))
	for i := 0; i < len(blobs); i++ {
		a[i] = make([]BLSFieldElement, len(blobs[i].Blob), len(blobs[i].Blob))
		for j := 0; j < len(blobs[i].Blob); j++ {
			a[i][j].SetBytesMod(blobs[i].Blob[j])
		}
	}
	return a
}

func BLSFieldElementsToBytes(bfes []BLSFieldElement) [][]byte {
	a := make([][]byte, len(bfes), len(bfes))
	for i := 0; i < len(bfes); i++ {
		copy(a[i], bfes[i].Bytes())
	}
	return a
}
