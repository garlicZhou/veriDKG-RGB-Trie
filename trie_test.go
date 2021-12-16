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

func TestStorage(t *testing.T) {
	db, err := leveldb.OpenFile("./db", nil)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	rgbtrie := new(db)

	fi, err := os.Open("c:/Users/85261/Desktop/1.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	i := 0
	for {
		if i > 100{
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
		//t.Triple.print()
		rgbtrie.tripleInsert(t)
		//fmt.Println("insert successful")
		fmt.Println(i)
		//rgbtrie.printTrie()

		//fmt.Println(string(a))
		//t := triple{subject: string(a)}
		//tp := tripleItem{Triple: &t, Address: string(a)}
		//hash := sha256.Sum256(a)
		//fmt.Println(hash)
		//fmt.Println(i)
		//db_hash := hash[:]
		//rgbtrie.Root.wordInsert(db_hash, tp,0, *db)
		//
		//x := i
		//
		//bytesBuffer := bytes.NewBuffer([]byte{})
		//binary.Write(bytesBuffer, binary.BigEndian, x)
		//db.Put(bytesBuffer.Bytes(), a,nil)
	}

	//rgbtrie.rootInsert()
	//rgbtrie_data,_ := rlp.EncodeToBytes(rgbtrie)
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
	fmt.Println("query result:", rgbtrie.searchTrie(strHash2))
	elapsed := time.Since(t1)
	fmt.Println("Query time:", elapsed.Nanoseconds())
	//db.Put(rgbtrie.RootHash, rgbtrie_data,nil)


	/*db.Put([]byte("key1"), []byte("value1"), nil)
	db.Put([]byte("key2"), []byte("value2"), nil)
	db.Put([]byte("key3"), []byte("value3"), nil)
	db.Put([]byte("key4"), []byte("value4"), nil)
	db.Put([]byte("key5"), []byte("value5"), nil)
	db.Put([]byte("food"), []byte("good"), nil)

	//循环遍历数据
	fmt.Println("循环遍历数据")
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Printf("key:%s, value:%s\n", iter.Key(), iter.Value())
	}
	iter.Release()*/

}
