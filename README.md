# go-dump
Quick hack for dumping memory stats to stderr for Go


Instrument code like so:
 ```
dump.Mem("some_func - BEFORE")
some_func()
dump.Mem("some_func - AFTER")
```

Output is like this: 
 ```
2019/10/10 14:37:48 some_func - BEFORE
2019/10/10 14:37:48   HeapAlloc  : 1.99G, delta: 408.46M
2019/10/10 14:37:48   HeapObjects: 18.36M, delta: 4.58M
2019/10/10 14:37:51 some_func - AFTER
2019/10/10 14:37:51   HeapAlloc  : 3.15G, delta: 1.15G
2019/10/10 14:37:51   HeapObjects: 22.98M, delta: 4.62M
```

Utility methods:

```
# Print number in human readable form
int b = 10
int k = 10000
int m = 10000000
int g = 10000000000

dump.Meg(b) // 10
dump.Meg(k) // 10K
dump.Meg(m) // 10M
dump.Meg(g) // 10G
```

