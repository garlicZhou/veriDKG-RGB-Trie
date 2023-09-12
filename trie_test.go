package RGBtrie

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	//數據庫操作
	db, err := leveldb.OpenFile("./db", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//創建Trie
	rgbtrie := newTrie(db)
	//三元組處理
	tripleNew := triple{subject: []byte("1"), predict: []byte("2"), object: []byte("3")}
	tripleNew.Hash()
	tp := tripleItem{Triple: tripleNew, Address: [32]byte{}}
	//三元組插入
	rgbtrie.tripleInsert(tp)
	//打印Trie
	rgbtrie.printTrie()
	//fmt.Println(rgbtrie.Root.hash)
	hash := [32]byte{200, 47, 93, 118, 40, 227, 190, 231, 249, 83, 101, 125, 11, 244, 189, 157, 254, 206, 4, 38, 3, 202, 233, 4, 176, 253, 16, 168, 109, 154, 96, 196}
	rgbtrie.Root.hash = hash
	//从数据库中重构RGBTrie
	rgbtrie.Root.reNewNode(db)
	rgbtrie.printTrie()

}

func TestStorage(t *testing.T) {
	db, err := leveldb.OpenFile("./db", nil)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	rgbtrie := newTrie(db)

	fi, err := os.Open("1.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	i := 0
	for {
		if i > 1000 {
			//rgbtrie.searchTrie()
			break
		}
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		i++
		stringA := string(a)
		item := strings.Fields(stringA)
		tripleNew := triple{subject: []byte(item[0]), predict: []byte(item[1]), object: []byte(item[2])}
		tripleNew.Hash()
		t := tripleItem{Triple: tripleNew, Address: [32]byte{}}
		fmt.Println(string(t.Triple.object))
		rgbtrie.tripleInsert(t)
		fmt.Println(i)
	}
	rgbtrie.printTrie()
	fmt.Println(rgbtrie.Root.child)
	fmt.Println(len(rgbtrie.Root.child))
	str := "http://dbpedia.org/ontology/Band"
	strHash := sha256.Sum256([]byte(str))
	strHash2 := strHash[:]
	t1 := time.Now()
	for i := 0; i < 2000; i++ {
		rgbtrie.searchTrie(strHash2)
	}
	//rgbtrie.printTrie()
	p, q := rgbtrie.searchTrie(strHash2)
	fmt.Println("query result:", p, q)
	elapsed := time.Since(t1)
	fmt.Println("Query time:", elapsed.Nanoseconds())

	fmt.Println("the result of veriy:", rgbtrie.verifyProof(p, q))

	fmt.Println("size of trie:", rgbtrie.getSizeOfTrie())

}
