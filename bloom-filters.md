 +++
author = ["Damian Gryski"]
date = "2014-12-01T08:00:00+00:00"
title = "Probabalistic Data Structures for Go: Bloom Filters"
series = ["Advent 2014"]
+++


Sometimes you have a lot of data.  If your problem can be solved by
[Hadoop](https://hadoop.apache.org/) and you happen to have a cluster lying
around, great!  However, sometimes you don't have enough memory or CPU to be
able to solve it exactly.  That's where probabilistic data structures come in.
These solutions let you trade-off accuracy in exchange for reduced memory
usage.

This article will briefly introduce one  common approximate data structures:
Bloom filters

## Bloom filters

A set is a collection of things.  You can add things to the set, and you can
query the set to see if an element has been added.  (We'll ignore deleting
elements from the set for now.) In Go, we might use a `map[string]struct{}` or
`map[string]bool` to represent a set.  If you query a map, you'll get back
"No, that element is not in the set." or "Yes, that element is in the set".

A Bloom filter is an 'approximate set'.  It supports the same two operations
(insert and query), but unlike a `map` the responses you'll get from a Bloom
filter query are "No, that element is not in the set" or "Yes, that element is
*probably* in the set".  How often it gives a false positive can be tuned by
the amount of space you want to use, but a good rule of thumb is that by
storing 10 *bits* per element, a Bloom filter will give a wrong answer about 1%
of the time.  One thing you can't do with a Bloom filter is iterate over it to
get back a list of items that have been inserted.

Under the hood, a Bloom filter is a bit-vector.  For each element you want to
put into the set, you hash it severeal times, and based on the value of each
hash you set certain bits in the bit-vector.  To query, you do the same hashing
and check the appropriate bits.  If any of the bits that are supposed to be set
are still 0, you know that element was never put into the set.  If they're all
ones, you only know that the element *might* be in the set.  Those bits could
have been set by hash collisions from other keys in the set.

## Applications of Bloom Filters

Here's an example where this approximate answer can still be useful.  Google
Chrome will warn you if you are about to visit a site Google has determined is
malicious.  How might we build this functionality into a web browser?  Chrome
could certainly query Google's servers for every URL, but that would slow down
our browsing since we now have to perform two network requests instead of just
one.  And since most URLs *aren't* malicious, the web service would spent most
of its time saying "Safe" to all the requests.

We could eliminate the network requests if Google Chrome had a local copy of
all the dangerous URLs that it could query instead.  But now instead of just
downloading a browser, we'd need to include a several gigabyte data file.

Lets see what happens if we put the malicious URLs into a Bloom filter instead.
First, unlike a Go map, Bloom filters use less space that the actual data they
are storing -- we no longer have to worry the huge download any more.  Now we
just have to check the Bloom filter before visiting a URL.  But what about the
wrong answers?  A Bloom filter of malicious URLs will never report a malicious
URL as "safe", it might only report a "safe" URL as "malicious".  For those cases,
false positives, we can still make the expensive call to Google's servers to
see if it really *is* malicious or one of the 1% false positives.

Bloom filters are also used in Cassandra and HBase as a way to avoid accessing
the disk searching for non-existent keys.  For more applications, I've listed
two papers under Further Reading.

## Bloom Filter Libraries

A standard Bloom filter is fairly easy to implement, so many people have.

Two popular ones are [willf/bloom](https://github.com/willf/bloom) and
[dataence/bloom](https://github.com/dataence/bloom).  This latter one
implements a number of different types of Bloom filters in addition to the
standard ones.

If you want to store lots of Bloom filters, there's also
[bloomd](https://github.com/armon/bloomd), a C network daemon for storing bloom
filters that also has a [Go bindings](https://github.com/geetarista/go-bloomd).
And the fine engineering team at [bitly](http://bit.ly) has written
[dablooms](https://github.com/bitly/dablooms), a high-performance Bloom filter
library in C that also has Go bindings.

TODO: add notes about other bloom filter types

## Further Reading*

* [Bloom Filters on Wikipedia](https://en.wikipedia.org/wiki/Bloom_filter)
* [Interactive Javascript Bloom filter demo](http://www.jasondavies.com/bloomfilter/)
* [Network Applications of Bloom Filters: A Survey](http://www.eecs.harvard.edu/~michaelm/NEWWORK/postscripts/BloomFilterSurvey.pdf)
* [Theory and Practice of Bloom Filters for Distributed Systems](http://www.dca.fee.unicamp.br/~chesteve/pubs/bloom-filter-ieee-survey-preprint.pdf)
