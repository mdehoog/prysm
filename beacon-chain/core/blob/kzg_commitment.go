package blob

type KZGCommitment [48]byte

func KZGsFromBytesArray(kzgs [][]byte) []KZGCommitment {
	a := make([]KZGCommitment, len(kzgs), len(kzgs))
	for i := 0; i < len(kzgs); i++ {
		copy(a[i][:], kzgs[i])
	}
	return a
}
