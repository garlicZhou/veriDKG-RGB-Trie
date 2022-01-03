package RGBtrie

import (
	"crypto/sha256"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
)

type node struct {
	parent    *node
	child     []*node
	childHash [][32]byte
	key       []byte
	value     []tripleItem
	hash      [32]byte
	isLeaf    bool
	isExtend  bool
    color     int8  //1:s, 2:p, 4:o, 3:sp, 5:so, 6:po, 7:spo
}

type nodekv struct {
	ChildHash [][32]byte
	Key       []byte
	Value     []tripleItem
	Hash      [32]byte
	IsLeaf    bool
	IsExtend  bool
	color     int8
}

type RGBtrie struct {
	Root *node
    RootHash [32]byte
	DB leveldb.DB
}

type proof struct {
	result []tripleItem
	merkleProof [][32]byte
}

func new(db *leveldb.DB) *RGBtrie {
    rgbtrie := &RGBtrie{
		Root: &node{isExtend: false, isLeaf: false, color: 0},
		DB: *db,
	}
	return rgbtrie
}

func (t *RGBtrie) putDb(db leveldb.DB) {
	t.DB = db
}

func (t *RGBtrie) PutRootHash() {
	if t.Root != nil  {
		t.RootHash = t.Root.hash
	}
}

func (t *RGBtrie) tripleInsert(item tripleItem)  {
	sub := item.Triple.subjectHash
	pre := item.Triple.predictHash
	obj := item.Triple.objectHash
	t.rootInsert(sub, item, 1)
	t.rootInsert(pre, item, 2)
	t.rootInsert(obj, item, 4)
}

func (t *RGBtrie) rootInsert(word []byte, item tripleItem, color int8) {
	flagRoot := true
	for _, j := range t.Root.child {
		if word[0] != j.key[0] {
			continue
		} else {
			j.nodeInsert(word, item, color, t.DB)
			flagRoot = false
			break
		}
	}
	if flagRoot {
		t.Root.child = append(t.Root.child, &node{key: word, parent: t.Root, value: []tripleItem{item}, isLeaf: true})
		t.Root.updateColor(color)
	}
	t.Root.updateColor(color)
	t.Root.updateHash(t.DB)
}

func (node1 *node) nodeInsert (word []byte, item tripleItem, color int8, db leveldb.DB) {
	lenOfWord := len(word)
	lenOfKey := len(node1.key)
	if lenOfKey == lenOfWord {
		flag := true
		for i, _ := range word {
			if word[i] == node1.key[i] {
				continue
			} else {
				node1.split(word, item, i - 1, color, db)
				flag = false
				break
			}
		}
		if flag {
			if node1.isLeaf == true || node1.isExtend == true {
				node1.value = append(node1.value, item)
				node1.updateColor(color)
				node1.updateHash(db)
			} else {
				node1.isExtend = true
				node1.value = append(node1.value, item)
				node1.updateColor(color)
				node1.updateHash(db)
			}
		}
	} else if lenOfKey > lenOfWord {
		flag := true
		for i, _ := range word {
			if word[i] == node1.key[i] {
				continue
			} else {
				node1.split(word, item, i - 1, color, db)
				flag = false
				break
			}
		}
		if flag {
			node1.split(word, item, lenOfWord - 1, color, db)
		}
	} else if lenOfKey < lenOfWord {
		flag := 0
		for i, _ := range node1.key {
			if word[i] == node1.key[i] {
				flag ++
				continue
			} else {
				node1.split(word, item, i - 1, color, db)
				flag = 0
				break
			}
		}
		if flag > 0 {
			flagRoot := true
			for _, j := range node1.child {
				if word[flag + 1] != j.key[0] {
					continue
				} else {
					j.nodeInsert(word[flag + 1 : lenOfWord], item, color, db)
					flagRoot = false
					break
				}
			}
			if flagRoot {
				node1.child = append(node1.child, &node{key: word[flag + 1 : lenOfWord], parent: node1, value: []tripleItem{item}})
			}
			node1.updateColor(color)
			node1.updateHash(db)
		} else {
			childNode := node{key: word[lenOfKey :], value: []tripleItem{item}, parent: node1, isLeaf: true, isExtend: false,color: color}
			childNode.updateHash(db)
			childNode.updateColor(color)
			node1.child = append(node1.child, &childNode)
			node1.updateHash(db)
			node1.updateColor(color)
		}
	}
}

func (node1 *node) split (word []byte, item tripleItem, lenOfSplit int, color int8, db leveldb.DB) {
	if lenOfSplit + 1 == len(word) {
		nodeNew := node{key: node1.key[lenOfSplit + 1 :], parent: node1, value: node1.value, isLeaf: node1.isLeaf,
			isExtend: node1.isExtend, color: node1.color, child: node1.child}
		node1.key = word
		node1.child = []*node{&nodeNew}
		node1.isLeaf = false
		node1.isExtend = true
		node1.value = []tripleItem{item}
	} else {
		nodeNew1 := node{key: node1.key[lenOfSplit + 1 :], parent: node1, value: node1.value, isLeaf: node1.isLeaf,
			isExtend: node1.isExtend, color: node1.color, child: node1.child}
		nodeNew2 := node{key: word[lenOfSplit + 1:], value: []tripleItem{item}, parent: node1, isLeaf: true,
			isExtend: false, color: color}
		node1.key = word[:lenOfSplit + 1]
		node1.child = []*node{&nodeNew1, &nodeNew2}
		node1.isLeaf = false
		node1.isExtend = false
		node1.value = nil
	}
	node1.updateColor(color)
	node1.updateHash(db)
}

func (node1 *node) updateHash(db leveldb.DB) {
	if len(node1.child) > 0 {
		node1.childHash = nil
		for _,j := range node1.child{
			node1.childHash = append(node1.childHash, j.hash)
		}
	} else {
		node1.childHash = nil
	}
	nodekv1 := nodekv{
		ChildHash: node1.childHash,
		Key:       node1.key,
		Value:     node1.value,
		Hash:      [32]byte{},
		IsLeaf:    node1.isLeaf,
		IsExtend:  node1.isExtend,
		color:     0,
	}
	if node1.isLeaf {
		data, _ := rlp.EncodeToBytes(nodekv1)
		node1.hash = sha256.Sum256(data)
		hash := node1.hash[:]
		db.Put(hash, data, nil)
		if node1.parent == nil {
			return
		} else {
			node1.parent.updateHash(db)
		}
	} else {
		data, _ := rlp.EncodeToBytes(nodekv1)
		for _, i := range node1.childHash {
			for _, j := range i{
				data = append(data, j)
			}
		}
		node1.hash = sha256.Sum256(data)
		hash := node1.hash[:]
		db.Put(hash, data, nil)
		if node1.parent == nil {
			return
		} else {
			node1.parent.updateHash(db)
		}
	}
}
//1:s, 2:p, 4:o, 3:sp, 5:so, 6:po, 7:spo
func (node1 *node) updateColor(color int8) {

		if (node1.color == 1 || node1.color == 2 || node1.color == 4) && node1.color != color {
			node1.color += color
		} else if (node1.color == 3 && color == 4) || (node1.color == 5 && color == 2) || (node1.color == 6 && color == 1){
			node1.color += color
		}
		if node1.parent != nil {
			node1.parent.updateColor(color)
		}
}

func (t *RGBtrie) printTrie() {
	fmt.Println("root: ")
	t.Root.printNode()
}

func (node1 *node) printNode() {
	if node1.parent != nil {
		fmt.Println("parents: ", node1.parent.hash)
	} else {
		fmt.Println("no parents")
	}

	//str := node1.hash[:]
	fmt.Print(" ", "keys:", node1.key, " ", "value:", node1.value, " ", "isLeaf: ",
		node1.isLeaf, " ", "isExtend: ", node1.isExtend, " ", "hash: ", node1.hash, " ")
	fmt.Println()
	if node1.isLeaf == false {
		for _,j := range node1.child {
			j.printNode()
		}
	}
}

func (t *RGBtrie) searchTrie(word []byte) []tripleItem {
	for _, j := range t.Root.child {
		if j.key[0] == word[0]{
           return j.searchNode(word)
		}
	}
	return []tripleItem{}
}

func (node1 *node) searchNode(word []byte, prf proof) (result proof) {
	if node1.isLeaf {
		if len(node1.key) != len(word) {
			return proof{}
		} else {
			for i, j := range node1.key {
				if word[i] == j {
					continue
				} else {
					return proof{}
				}
			}
			prf.result = node1.value
			return prf
		}
	} else if node1.isExtend {
		if len(node1.key) == len(word) {
			for i, j := range node1.key {
				if word[i] == j {
					continue
				} else {
					return []tripleItem{}
				}
			}
			return node1.value
		} else if len(node1.key) < len(word) {
			for i, j := range node1.key {
				if word[i] == j {
					continue
				} else {
					return []tripleItem{}
				}
			}
			for _, q := range node1.child {
				if word[len(node1.key)] == q.key[0]{
					return q.searchNode(word[len(node1.key) :])
				}
			}
			return []tripleItem{}
		}
	} else {
		if len(node1.key) < len(word) {
			for i, j := range node1.key {
				if word[i] == j {
					continue
				} else {
					return []tripleItem{}
				}
			}
			for _, q := range node1.child {
				if word[len(node1.key)] == q.key[0]{
					return q.searchNode(word[len(node1.key) :])
				}
			}
			return []tripleItem{}
		}
	}
	return []tripleItem{}
}

func (t *RGBtrie) getProof() {

}

func (t *RGBtrie) verifyProof() {

}


