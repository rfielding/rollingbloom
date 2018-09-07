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
