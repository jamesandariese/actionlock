# `ActionLock`

_A lock that does something just before the value changes._

## Description

`ActionLock` orders access to a resource in contention so that goroutines
depending on a specific global value being set can run in parallel while the
same global value is required.

`ActionLock`'s purpose in this is to further support use of the global resource
by allowing it to be manipulated any time the ActionLock value changes.

Use it to control the active network, frequency/band, or other physical or
limited resource.

## Operation

Create an `ActionLock` with a callback function.  Whenever the value in the
`ActionLock` changes, the callback will be called.  Use this callback function
to do things like change active band on a radio or a network VLAN.

In each code path which requires coordination of that value, lock the value
with `LockValue` as normal other than providing the needed value.  Ensure you
`UnlockValue` as well.  Using `defer` for this is best practice.

Also see the example in `actionlock_test.go`