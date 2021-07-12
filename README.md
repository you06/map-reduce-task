# Map Reduce Task

[中文版](./README-zh.md)

In this task, you need to complete a map-reduce workflow and simulate it in a single process.

Suppose we have an extremely large file with words seperated by space in it, e.g., 100GB, which is not able to loaded into RAM memory. Your task is to find the first non-repeat word. There are basic following rules.

- Only allow scan the text file once
- Use less IO operation as possible
- RAM limitation is 16G
- The implementation must be stable, in another word, even missing the first non-repeat word in a very low probability is not allowed.

There is a basic framework and a test data generator. In your implementation, there are some additional rules.

- Your code which read files must be under `datanode` directory
- Your code which collect the middle result or final result must be placed under `calcnode` directory
- Each node must run in individual [goroutine](https://golang.org/doc/effective_go#goroutines)
- Data exchange between node instances must be through [channels](https://golang.org/doc/effective_go#channels)

It would be bonus if there is a data exchange format which supports serialization and deserialization in your implementation.

The RAM limitation is for single node, you don't need to care the total memory usage when simulate in a single machine.
