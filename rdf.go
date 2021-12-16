package RGBtrie

import (
	"crypto/sha256"
	"fmt"
)

type triple struct {
	subject      []byte
	subjectHash  []byte
	predict      []byte
	predictHash  []byte
	object       []byte
	objectHash   []byte
}

type tripleItem struct {
	Triple  triple
	Address [32]byte
}

type rdfFragment struct {
	Triples []tripleItem
}

func (t *triple) Hash()  {
	subjectHash := sha256.Sum256(t.subject)
	t.subjectHash = subjectHash[:]
	predictHash := sha256.Sum256(t.predict)
	t.predictHash = predictHash[:]
	objectHash := sha256.Sum256(t.object)
	t.objectHash = objectHash[:]
}

func (t *triple) print() {
	fmt.Println(string(t.subject), string(t.predict), string(t.object))
}

/*func readRDF() *rdfFragment {
	return
}*/