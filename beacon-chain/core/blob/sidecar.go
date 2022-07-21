package blob

import (
	kbls "github.com/kilic/bls12-381"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/consensus-types/interfaces"
	types "github.com/prysmaticlabs/prysm/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	v1 "github.com/prysmaticlabs/prysm/proto/engine/v1"
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"math/big"
)

var (
	ErrInvalidBlobSlot            = errors.New("invalid blob slot")
	ErrInvalidBlobBeaconBlockRoot = errors.New("invalid blob beacon block root")
	ErrInvalidBlobsLength         = errors.New("invalid blobs length")
	ErrCouldNotComputeCommitment  = errors.New("could not compute commitment")
	ErrMissmatchKzgs              = errors.New("missmatch kzgs")
)

// VerifyBlobsSidecar verifies the integrity of a sidecar.
// def verify_blobs_sidecar(slot: Slot, beacon_block_root: Root,
//                         expected_kzgs: Sequence[KZGCommitment], blobs_sidecar: BlobsSidecar):
//    assert slot == blobs_sidecar.beacon_block_slot
//    assert beacon_block_root == blobs_sidecar.beacon_block_root
//    blobs = blobs_sidecar.blobs
//    assert len(expected_kzgs) == len(blobs)
//    for kzg, blob in zip(expected_kzgs, blobs):
//        assert blob_to_kzg(blob) == kzg
func VerifyBlobsSidecar(slot types.Slot, beaconBlockRoot [32]byte, expectedKZGs [][]byte, blobsSidecar *eth.BlobsSidecar) error {
	if slot != blobsSidecar.BeaconBlockSlot {
		return ErrInvalidBlobSlot
	}
	if beaconBlockRoot != bytesutil.ToBytes32(blobsSidecar.BeaconBlockRoot) {
		return ErrInvalidBlobBeaconBlockRoot
	}
	if len(expectedKZGs) != len(blobsSidecar.Blobs) {
		return ErrInvalidBlobsLength
	}

	aggregatedPoly, aggregatedPolyCommitment, err := computeAggregatedPolyAndCommitment(blobsSidecar.Blobs, expectedKZGs)
	if err != nil {
		return err
	}

	x, err := hashToBLSField2(aggregatedPoly, aggregatedPolyCommitment)
	if err != nil {
		return err
	}

	y, err := evaluatePolynomialInEvaluationForm(aggregatedPoly, x)

	return verifyKZGProof(aggregatedPolyCommitment, x, y, blobsSidecar.AggregatedProof)
}

func computeAggregatedPolyAndCommitment(blobs []*v1.Blob, kzgCommitments [][]byte) ([]BLSFieldElement, KZGCommitment, error) {
	r, err := hashToBLSField(blobs, kzgCommitments)
	if err != nil {
		return nil, KZGCommitment{}, err
	}
	rPowers := computePowers(r, len(kzgCommitments))

	blobsBLS := BlobsToBLSFieldElements(blobs)
	aggregatedPoly, err := matrixLincomb(blobsBLS, rPowers)
	if err != nil {
		return nil, KZGCommitment{}, err
	}

	kzgs := KZGsFromBytesArray(kzgCommitments)
	aggregatedPolyCommitment, err := lincomb(kzgs, rPowers)
	if err != nil {
		return nil, KZGCommitment{}, err
	}

	return aggregatedPoly, aggregatedPolyCommitment, nil
}

func hashToBLSField(blobs []*v1.Blob, expectedKZGs [][]byte) (BLSFieldElement, error) {
	bwk := eth.BlobsWithKzgs{
		Blobs: blobs,
		Kzgs:  expectedKZGs,
	}
	htr, err := bwk.HashTreeRoot()
	if err != nil {
		return BLSFieldElement{}, err
	}
	i := BLSFieldElement{}
	i.SetBytesMod(htr[:])
	return i, nil
}

func hashToBLSField2(poly []BLSFieldElement, commitment KZGCommitment) (BLSFieldElement, error) {
	ap := eth.AggregatedPoly{
		AggregatedPoly:           BLSFieldElementsToBytes(poly),
		AggregatedPolyCommitment: commitment[:],
	}
	htr, err := ap.HashTreeRoot()
	if err != nil {
		return BLSFieldElement{}, err
	}
	i := BLSFieldElement{}
	i.SetBytesMod(htr[:])
	return i, nil
}

func computePowers(x BLSFieldElement, n int) []BLSFieldElement {
	currentPower := BLSFieldElement{}
	currentPower.SetFromBigMod(big.NewInt(1))
	var powers []BLSFieldElement
	for i := 0; i < n; i++ {
		powers = append(powers, currentPower)
		currentPower.MulMod(&x)
	}
	return powers
}

func lincomb(vectors []KZGCommitment, scalars []BLSFieldElement) (KZGCommitment, error) {
	var err error
	points := make([]*kbls.PointG1, len(vectors), len(vectors))
	for i := 0; i < len(vectors); i++ {
		if points[i], err = kbls.NewG1().FromCompressed(vectors[i][:]); err != nil {
			return KZGCommitment{}, err
		}
	}
	scalarsFr := make([]*kbls.Fr, len(scalars), len(scalars))
	for i := 0; i < len(scalars); i++ {
		scalarsFr[i] = (*kbls.Fr)(&scalars[i])
	}

	rg := kbls.PointG1{}
	if _, err = kbls.NewG1().MultiExp(&rg, points, scalarsFr); err != nil {
		return KZGCommitment{}, err
	}

	r := KZGCommitment{}
	copy(r[:], kbls.NewG1().ToCompressed(&rg))
	return r, nil
}

func matrixLincomb(vectors [][]BLSFieldElement, scalars []BLSFieldElement) ([]BLSFieldElement, error) {
	if len(vectors) != len(scalars) {
		return nil, errors.New("vectors and scalar vectors should be the same length")
	}
	r := make([]BLSFieldElement, len(vectors[0]), len(vectors[0]))
	for i := 0; i < len(vectors[0]); i++ {
		a := scalars[i]
		for j, x := range vectors[i] {
			r[j].AddMod(x.MulMod(&a))
		}
	}
	return r, nil
}

func evaluatePolynomialInEvaluationForm(poly []BLSFieldElement, x BLSFieldElement) (BLSFieldElement, error) {
	// TODO
}

func verifyKZGProof(commitment KZGCommitment, x BLSFieldElement, y BLSFieldElement, proof []byte) error {
	// TODO
}

func BlockContainsKZGs(b interfaces.BeaconBlock) bool {
	if blocks.IsPreEIP4844Version(b.Version()) {
		return false
	}
	blobKzgs, _ := b.Body().BlobKzgs()
	return len(blobKzgs) != 0
}
