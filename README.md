# Critbit

Critbit is a Go package which implements a critbit tree.
A critbit tree is useful for quickly finding a string within
a tree of strings. It naturally stores strings in sorted order,
so they can be also trivially be retrieved from the tree in
sorted order.

The implementation is novel in that it uses two arrays to
maintain all nodes, instead of relying on pointers. Aspects
of this implementation were influenced by:

https://github.com/mb0/critbit

and

https://github.com/glk/critbit

## Example
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

