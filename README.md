# RGB-Trie in VeriDKG

This project is a go-language implementation of VeriDKG's core module, RGB-Trie.

To test the performance of RGB-Trie in real blockchain scenarios, please replace the trie.go file in the ETHMST project with the trie.go file in this project for testing, the ETHMST repository address is: [](https://github.com/garlicZhou/ETHMS)

Task: 测试加入RGB-Trie对以太坊的成本增加

1. 只测试吞吐量影响，不测试存储开销

2. Block header中建一个rgbroot 存放rgbtrie的根哈希，在以太坊里新建一个目录把rdf.go, trie.go放进去

3. 在以太坊生成区块时，创建一个rgb-trie

   ```
   rgbtrie := newTrie(db)
   ```

   判断上一个区块header中的rgbroot是否为空，如果不为空

   ```
   rgbtrie.Root.reNewNode(db)
   ```

   对每条教育中的三元组进行处理，转为三个哈希

   ```
   tripleNew := triple{subject: []byte("1"), predict: []byte("2"), object: []byte("3")}
   tripleNew.Hash()
   tp := tripleItem{Triple: tripleNew, Address: [32]byte{}}
   ```

​       将每条交易中的三元组插入rgb-trie

```
 rgbtrie.tripleInsert(tp)
```

​       将rgbtrie的root hash存入block header

```
block.rgbroot = rgbtrie.Root.hash
```

4.如果以太坊不兼容LevelDB， 就改成与以太坊一致的ethdb.Database
