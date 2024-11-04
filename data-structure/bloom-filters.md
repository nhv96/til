# Bloom Filter
Bloom filter is a probabilistic data structure used to quickly determine if an element is in a set.

This data structure is very space efficient.

It return either "probabilty the item is in the list" or "the item is definitely not in the list".

Due to it's probabilistic, it trades off the accuracy for memory consume.

Use for cases where we can tolerate for false positive, but not false negative.

# How it works
Bloom filter is a bit array of *m* bits, each set to 0.

We also need *k* hash functions, outputs from these hash functions will map to set positions of the bit array. Hash functions must be randomly distributed.

When add an element, feed it to the *k* hash functions to get *k* array positions, then set these positions in the bit array to 1.

To test if an element is in the array, feed it to the *k* hash functions and get the *k* array positions, then check if these positions in the array all containing value 1. If there is value 0, the element is not in the array.