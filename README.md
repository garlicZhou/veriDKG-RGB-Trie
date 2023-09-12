# RGB-Trie in VeriDKG

This project is a go-language implementation of VeriDKG's core module, RGB-Trie.

To test the performance of RGB-Trie in real blockchain scenarios, please replace the trie.go file in the ETHMST project with the trie.go file in this project for testing, the ETHMST repository address is: 
[ETHMST](https://github.com/garlicZhou/ETHMST)

Task: Test the additional cost of adding RGB-Trie to Ethereum.

1. Create an empty hash in the Block header to store the root hash of rgbtrie, and create a new directory in Ethereum to include rdf.go and trie.go.

2. Create an RGB-Trie when generating blocks in Ethereum

   ```
   rgbtrie := newTrie(db)
   ```

   Determine whether the rgbroot in the previous block header is empty, and if it is not

   ```
   rgbtrie.Root.reNewNode(db)
   ```

   Process the triples in each transaction and convert them into three hashes

   ```
   tripleNew := triple{subject: []byte("1"), predict: []byte("2"), object: []byte("3")}
   tripleNew.Hash()
   tp := tripleItem{Triple: tripleNew, Address: [32]byte{}}
   ```

   Insert triples from each transaction into the RGB-Trie

   ```
   rgbtrie.tripleInsert(tp)
   ```

​   Save the root hash of the RGB-Trie into the block header

   ```
   block.rgbroot = rgbtrie.Root.hash
   ```

3.If Ethereum is not compatible with LevelDB, change to ETHDB.Database that is consistent with Ethereum
