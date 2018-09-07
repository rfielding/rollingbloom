Rolling Bloom Filter
====================

Adding entries into a bloom filter is cumulative.
There isn't a simple way to remove items from it.
Instead, a simple strategy of triple-buffering the filter can be used.

Every Add and Test operation passes in a logical clock.
The filter has an interval over which it circulates through the buffers.
An Add will write the entry into current and next filter.
A test will read the current filter.
When either Add or Test are called, we may clear the current filter and move forward to the next one.

Uses
====

In Caching, if we can get the hash collision rate low enough, we may be able to ignore it for some applications.
In particular, not allowing too much stuff to be written into a bloom filter keeps its accuracy high.
We also may be dealing with data that expires.

- Suppose that you have to calculate a yes/no answer to an expensive deterministic question (because you need to get the answer from a remote service, very repetitively).  You can cache the answer along with the parameters that went into it.

You can fill the cache full of positive and negative results:
```
yes = rfielding canRead dev google
no = rfielding canRead ceo google
allow IAM {"name":"rob.fielding@gmail.com", "groups":["dev","rnd"], "expires-month":"2018-September"}
as {"name":"rob.fielding"} read {"label":"topstupid"}
```

Since the whole point of a bloom filter is to store a very large number of items,
the fact that there are many combinations is less of a problem.

> What is most useful about it is that the actual size of the things being hashed doesn't really matter, as only the hash is stored.

Clever encoding of potentially repetitive questions lets us accumulate a cache of assertions.
Making it rolling deals with the fact that the answer may only be valid for a limited time.
> TODO: taking the round into account in the hashes would allow us to more gracefully deal with false positives, as it may be ok to possibly get a wrong answer with low probability, and as long as discovering this problem does not persist.  (ie: including the round in the filter hash)

