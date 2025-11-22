# Critbit

Critbit is a Go package which implements a critbit tree.
A critbit tree is useful for quickly finding a string within
a tree of strings. It naturally stores strings in sorted order,
so they can be also trivially be retrieved from the tree in
sorted order.

See the on-line Go doc for this package at:

https://godoc.org/github.com/gilramir/critbit

The implementation is novel in that it uses two arrays to
maintain all nodes, instead of relying on pointers. Aspects
of this implementation were influenced by:

https://github.com/mb0/critbit

and

https://github.com/glk/critbit

## Example
```
    package main

    import (
        "fmt"
        "github.com/gilramir/critbit"
    )

    func main() {
        tree := critbit.New(0)
        ok, err = tree.Insert("gamma", 300)
        ok, err = tree.Insert("beta", 200)
        ok, err = tree.Insert("alpha", 100)

        for kv := range tree.GetKeyValueTuples() {
            fmt.Printf("%s = %d\n", kv.Key, kv.Value.(int))
        }
    }
```

## Methods
* **Delete** - delete a key
* **Dump** - print the trie's representation to stdout, for debugging
* **Get** - get a key's value
* **GetHasPrefix** - find the first key that starts with a prefix,
    and return the KeyValueTuple
* **GetKeyValueTuples** - get all key/value tuples
* **GetKeyValueTupleChan** - get a channel to read all key/value tuples
* **Insert** - insert a new key/value, without updating an existing key
* **Keys** - get all keys
* **Length** - get the number of keys
* **Louds()** - get the LOUDS representation of the trie
* **MemorySizeBytes** - get an approximation of how much memory the trie is
  using
* **SaveDot** - output the tree in graphviz/dot format
* **Split** - split a trie into 2 even tries
* **SplitAt** - split a trie into 2 tries at any key
* **Update** - update an existing key's value, without inserting a new key
