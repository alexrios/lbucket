/*
Package lbucket provides support for using leaky buckets on your app.

Creating a new Leaky bucket informing 3 as the bucket capacity and the frequency how the bucket leaks.
	 NewTickLeakyBucket(3, 1 * time.Second)

Calling Refill() you're add more volume in the bucket.
	bucket := NewTickLeakyBucket(3, 1 * time.Second)
	err := bucket.Refill()

Note: When the bucket capacity is reached a ErrBucketReachedCap will be returned until the bucket leaks once again.
A simple errors.Is could be used in this scenario:
	errors.Is(err, ErrBucketReachedCap)

To stop the bucket leaking you can call the Fix() method.
	...
	bucket.Fix()
	...

Use Size() to get to current volume in the bucket.
	...
	bucket.Size()
	...
*/
package lbucket
