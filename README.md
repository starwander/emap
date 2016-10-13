## Enhanced Golang Map
[![Build Status](https://travis-ci.org/EthanZhuang/EMap.svg?branch=master)](https://travis-ci.org/EthanZhuang/EMap)
[![codecov](https://codecov.io/gh/EthanZhuang/EMap/branch/master/graph/badge.svg)](https://codecov.io/gh/EthanZhuang/EMap)
[![Go Report Card](https://goreportcard.com/badge/github.com/EthanZhuang/EMap)](https://goreportcard.com/report/github.com/EthanZhuang/EMap)
[![GoDoc](https://godoc.org/github.com/EthanZhuang/EMap?status.svg)](https://godoc.org/github.com/EthanZhuang/EMap)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

EMap implements an enhanced map in [Golang](https://golang.org/).
The main enhancement of emap over original golang map is the support of searching index.
* Values in the emap can have one or more indices which can be used to search or delete.
* Key in the emap must be unique as same as the key in the golang map.
* Index in the emap is an N to M relation which mean a value can have multi indices and multi values can have one same index.

##Several Choice
#####Generic EMap
* The generic emap has no restrict for the type of its key, value and index.
* The generic emap has a read-write locker inside so it is concurrent safe.

#####Expirable EMap
* The expirable emap has no restrict for the type of its key, value and index.
* The expirable emap has a read-write locker inside so it is concurrent safe.
* The expirable emap will check all the values in the emap with the period of input interval(milliseconds). If a value is expired, it will be deleted automatically.

#####Strict EMap
* The types of key, value and index used in the strict emap are determined during initialization by the sample inputs.
* All methods of the strict emap must use the same type of the sample inputs otherwise an error will be returned.
* The strict emap has a read-write locker inside so it is concurrent safe.

#####Unlock EMap
* The unlock emap has no restrict for the type of its key, value and index.
* The unlock emap has no locker or mutex inside, so it is not concurrent safe.
* It is only suitable for those models like Event Loop to achieve better performance.

##Requirements
#####Download this package

    go get github.com/EthanZhuang/EMap

#####Implements ExpirableValue interface of this package for all values if ExpirableEmap is chosen
```go
// ExpirableValue is the interface which must be implemented by all the value in the expirable EMap.
type ExpirableValue interface {
	// IsExpired returns if the value is expired.
	// If true, the value will be deleted automatically.
	IsExpired() bool
}
```
## Basic Operations
* Insert: pushes a new value into emap with input key and indices.
* FetchByKey: gets the value in the emap by input key.
* FetchByIndex: gets the all values in the emap by input index.
* DeleteByKey: deletes the value in the emap by input key.
* DeleteByIndex: deletes all the values in the emap by input index.
* AddIndex: add the input index to the value in the emap of the input key.
* RemoveIndex: remove the input index from the value in the emap of the input key.
* KeyNum: returns the total key number in the emap.
* KeyNumOfIndex: returns the total key number of the input index in the emap.
* IndexNum: returns the total index number in the emap.
* IndexNumOfKey: returns the total index number of the input key in the emap.
* HasKey: returns if the input key exists in the emap.
* HasIndex: returns if the input index exists in the emap.

## Higher-order Operations
* Transform:
- Transform is a higher-order operation which apply the input callback function to each key-value pair in the emap.
- Any error returned by the callback function will interrupt the transforming and the error will be returned.
- If transform successfully, a new golang map is created with each key-value pair returned by the input callback function.

* Foreach:
- Foreach is a higher-order operation which apply the input callback function to each key-value pair in the emap.
- Since the callback function has no return, the foreach procedure will never be interrupted.
- A typical usage of Foreach is apply a closure.


## Example

```go

```

## Reference

[GoDoc](https://godoc.org/github.com/EthanZhuang/EMap)

## LICENSE

EMap source code is licensed under the [Apache Licence, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).
