# 0.1.0

* initial version, all tests passed

# 1.0.0

* public release
* add: new interface Repeater
* add: new concept - OpWrapper
* add: new function to create Repeater with wrapping operations
* add: new WithContext - a repeater that checks for context errors before operation call

# 1.1

* add: IsTemporary and IsStop to errors
* add: FnOnError to operations
* add: missing tests to operations

# 1.2

* add: Nope to operations
* add: new function to create Repeater with constructor and destructor - Cpp

# 1.3

* add: Once and FnRepeat to Repeater and global functions

# 1.4

* fix: Once calls global compose
* ref: some changes in Cpp concept. It is transparent for input errors now, it also panics if D fails
* add: Done, FnDone, FnOnlyOnce
* add: 100% test coverage
