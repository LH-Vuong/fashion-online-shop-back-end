package key

import (
	cid "github.com/ipfs/go-cid"
)

func GetHashFromCidString(c string) (string, error) {
	_cid, err := cid.Parse(c)
	if err != nil {
		return "", err
	}
	mh := _cid.Hash()
	return mh.HexString(), err
}
