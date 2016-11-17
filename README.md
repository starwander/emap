## Enhanced Golang Map
[![Build Status](https://travis-ci.org/starwander/EMap.svg?branch=master)](https://travis-ci.org/starwander/EMap)
[![codecov](https://codecov.io/gh/starwander/EMap/branch/master/graph/badge.svg)](https://codecov.io/gh/starwander/EMap)
[![Go Report Card](https://goreportcard.com/badge/github.com/starwander/EMap)](https://goreportcard.com/report/github.com/starwander/EMap)
[![GoDoc](https://godoc.org/github.com/starwander/EMap?status.svg)](https://godoc.org/github.com/starwander/EMap)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

EMap implements an enhanced map in [Golang](https://golang.org/).
The main enhancement of emap over original golang map is the support of searching index.
* Values in the emap can have one or more indices which can be used to search or delete.
* Key in the emap must be unique as same as the key in the golang map.
* Index in the emap is an N to M relation which mean a value can have multi indices and multi values can have one same index.

##Several Choices
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

    go get github.com/starwander/EMap

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
package main

import (
	"fmt"
	"github.com/starwander/EMap"
	"time"
)

type employee struct {
	name       string
	gender     string
	department string
	retired    bool
	//hobby      []string
}

type emplyeeByDept struct {
	name       string
	department string
}

func (e *employee) IsExpired() bool {
	return e.retired
}

func main() {
	emap := emap.NewExpirableEMap(1000)

	emap.Insert("Tom", &employee{"Tom", "Male", "R&D", false}, "R&D", "swimming", "hiking")
	emap.Insert("Jerry", &employee{"Jerry", "Male", "Sales", false}, "Sales", "football", "hiking")
	emap.Insert("Carol", &employee{"Carol", "Female", "HR", false}, "HR", "movie")
	emap.Insert("Jessie", &employee{"Jessie", "Female", "Sales", false}, "Sales", "hiking", "cycling")

	fmt.Println(emap.KeyNum())                //4
	fmt.Println(emap.IndexNum())              //8
	fmt.Println(emap.IndexNumOfKey("Tom"))    //3
	fmt.Println(emap.KeyNumOfIndex("hiking")) //3

	fmt.Println(emap.FetchByKey("Carol")) // &{Carol Female HR false} <nil>
	allSales, _ := emap.FetchByIndex("Sales")
	for _, each := range allSales {
		//&{Jerry Male Sales false}
		//&{Jessie Female Sales false}
		fmt.Println(each.(*employee))
	}

	emap.RemoveIndex("Jerry", "hiking")
	allHiking, _ := emap.FetchByIndex("hiking")
	for _, each := range allHiking {
		//&{Tom Male R&D false}
		//&{Jessie Female Sales false}
		fmt.Println(each.(*employee))
	}
	emap.AddIndex("Jessie", "movie")
	allMovie, _ := emap.FetchByIndex("movie")
	for _, each := range allMovie {
		//&{Carol Female HR false}
		//&{Jessie Female Sales false}
		fmt.Println(each.(*employee))
	}

	emap.DeleteByKey("Jerry")
	fmt.Println(emap.KeyNum()) //3

	nameAndDept, _ := emap.Transform(func(n interface{}, e interface{}) (interface{}, error) {
		return &emplyeeByDept{n.(string), e.(*employee).department}, nil
	})
	for _, each := range nameAndDept {
		//&{Jessie Sales}
		//&{Tom R&D}
		//&{Carol HR}
		fmt.Println(each.(*emplyeeByDept))
	}

	var totalMaleNum, totalFemalNum int
	emap.Foreach(func(n interface{}, e interface{}) {
		if e.(*employee).gender == "Male" {
			totalMaleNum++
		} else {
			totalFemalNum++
		}
	})
	fmt.Println("Total male number:", totalMaleNum, "Total female number:", totalFemalNum) //Total male number: 1 Total female number: 2

	carol, _ := emap.FetchByKey("Carol")
	carol.(*employee).retired = true
	time.Sleep(1100 * time.Millisecond)
	fmt.Println(emap.FetchByKey("Carol")) //<nil> key not exist
}
```

## Reference

[GoDoc](https://godoc.org/github.com/starwander/EMap)

## LICENSE

EMap source code is licensed under the [Apache Licence, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).
