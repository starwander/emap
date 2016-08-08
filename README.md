# Enhanced Golang Map: support one unique key and multi search indices for each value #
Add some multi-index support into the original golang map:

1. Each value in the emap must has one unique key.

2. Each value in the emap can have multi indices.

3. One index can related to multi values in the emap.

## Interfaces ##
EMap:
- Insert(key interface{}, value interface{}, indices ...interface{}) error
- FetchByKey(key interface{}) (interface{}, error)
- FetchByIndex(index interface{}) ([]interface{}, error)
- DeleteByKey(key interface{}) error
- DeleteByIndex(index interface{}) error
- AddIndex(key interface{}, index interface{}) error
- RemoveIndex(key interface{}, index interface{}) error
- KeyNum() int
- KeyNumOfIndex(index interface{}) int
- IndexNum() int
- IndexNumOfKey(key interface{}) int
- HasKey(key interface{}) bool
- HasIndex(index interface{}) bool
- Transform(callback func(interface{}, interface{})(interface{}, error)) (map[interface{}]interface{}, error)
- Foreach(callback func(interface{}, interface{}) error) error

ExpirableValue:
- IsExpired() bool

## Several Implementations of Emap##
- generic_emap: The basic implementation of emap. The key, index and value can be anything.
- strict_emap: Add some type check into all interfaces. The type of key, index and value is appointed during initialization. Use different types  later should fail.
- expirable_emap: Emap will check each value for expiration every interval appointed during initialization. Value added into expirable emap must implement ExpirableValue interface.
- unlock_emap: Emap will not lock anything so it's not concurrent safe. This is only suitable for those Event Loop code who can use unlock emap to achieve better performance.

## Example ##
EMap is quite easy to use. Check the tests for more details.

## License ##
This library is under the [MIT License](http://opensource.org/licenses/MIT)
