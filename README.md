# lbucket
[![Go Reference](https://pkg.go.dev/badge/github.com/alexrios/lbucket.svg)](https://pkg.go.dev/github.com/alexrios/lbucket?tab=doc)

lbucket is an idiomatic Go leaky bucket implementation.

The library make use of plain old Go stdlib; in other words, there are no third-party dependencies.

### How to use
Package lbucket provides support for using leaky buckets on your app.

Creating a new Leaky bucket informing 3 as the bucket capacity and the frequency how the bucket leaks.
```go
NewTickLeakyBucket(3, 1 * time.Second)
```

Calling `Refill()` you're add more volume in the bucket.
```go
bucket := NewTickLeakyBucket(3, 1 * time.Second)
err := bucket.Refill()
```
Note: When the bucket capacity is reached a `ErrBucketReachedCap` will be returned until the bucket leaks once again.
A simple `errors.Is` could be used in this scenario:
```go
errors.Is(err, ErrBucketReachedCap)
```

To stop the bucket leaking you can call the `Fix()` method.
```go
bucket.Fix()
```

Use Size() to get to current volume in the bucket.
```go
bucket.Size()
```