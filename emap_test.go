// Copyright(c) 2016 Ethan Zhuang <zhuangwj@gmail.com>.

package emap

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"time"
)

type EMap interface {
	Insert(key interface{}, value interface{}, indices ...interface{}) error
	FetchByKey(key interface{}) (interface{}, error)
	FetchByIndex(index interface{}) ([]interface{}, error)
	DeleteByKey(key interface{}) error
	DeleteByIndex(index interface{}) error
	AddIndex(key interface{}, index interface{}) error
	RemoveIndex(key interface{}, index interface{}) error
	KeyNum() int
	KeyNumOfIndex(index interface{}) int
	IndexNum() int
	IndexNumOfKey(key interface{}) int
	HasKey(key interface{}) bool
	HasIndex(index interface{}) bool
	Transform(callback func(interface{}, interface{}) (interface{}, error)) (map[interface{}]interface{}, error)
	Foreach(callback func(interface{}, interface{}))
}

var _ = Describe("Tests of emap", func() {
	Context("one unique key and multi indices", func() {
		DescribeTable("Given an empty emap, when add a new item, it should be able to get by key or index later.", func(emap EMap) {
			Expect(emap.HasKey("key1")).To(Equal(false))
			err := emap.Insert("key1", "value1", "index1", "index2", "index3")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(emap.HasKey("key1")).To(Equal(true))
			result1, err := emap.FetchByKey("key1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result1).To(BeEquivalentTo("value1"))

			result2, err := emap.FetchByIndex("index1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result2).To(BeEquivalentTo([]interface{}{"value1"}))

			result3, err := emap.FetchByIndex("index2")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result3).To(BeEquivalentTo([]interface{}{"value1"}))

			result4, err := emap.FetchByIndex("index3")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result4).To(BeEquivalentTo([]interface{}{"value1"}))
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an emap with key1, when add a new item with the same key1, it should fail.", func(emap EMap) {
			err := emap.Insert("key1", "value1", "index1", "index2")
			Expect(err).ShouldNot(HaveOccurred())

			err = emap.Insert("key1", "value2", "index3")
			Expect(err).Should(HaveOccurred())

			Expect(emap.KeyNum()).To(BeEquivalentTo(1))
			Expect(emap.IndexNum()).To(BeEquivalentTo(2))
			Expect(emap.IndexNumOfKey("key1")).To(BeEquivalentTo(2))
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an empty emap, when delete a item with the key1, it should fail.", func(emap EMap) {
			err := emap.DeleteByKey("key1")
			Expect(err).Should(HaveOccurred())
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an emap with multi values, when delete by key, it should delete the value and indices of the key.", func(emap EMap) {
			err := emap.Insert("key1", "value1", "index1", "index2")
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.Insert("key2", "value2", "index3")
			Expect(err).ShouldNot(HaveOccurred())
			result1, err := emap.FetchByIndex("index2")
			Expect(result1).To(BeEquivalentTo([]interface{}{"value1"}))

			Expect(emap.KeyNum()).To(BeEquivalentTo(2))
			Expect(emap.IndexNum()).To(BeEquivalentTo(3))
			Expect(emap.IndexNumOfKey("key1")).To(BeEquivalentTo(2))
			Expect(emap.IndexNumOfKey("key2")).To(BeEquivalentTo(1))

			err = emap.DeleteByKey("key1")
			Expect(err).ShouldNot(HaveOccurred())
			_, err = emap.FetchByKey("key1")
			Expect(err).Should(HaveOccurred())
			result2, err := emap.FetchByKey("key2")
			Expect(result2).To(BeEquivalentTo("value2"))
			Expect(emap.KeyNum()).To(BeEquivalentTo(1))
			Expect(emap.IndexNum()).To(BeEquivalentTo(1))
			Expect(emap.IndexNumOfKey("key1")).To(BeEquivalentTo(0))
			Expect(emap.IndexNumOfKey("key2")).To(BeEquivalentTo(1))
			Expect(emap.HasKey("key1")).To(Equal(false))
			Expect(emap.HasKey("key2")).To(Equal(true))

			err = emap.DeleteByKey("key2")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(emap.HasKey("key1")).To(Equal(false))
			Expect(emap.KeyNum()).To(BeEquivalentTo(0))
			Expect(emap.IndexNum()).To(BeEquivalentTo(0))
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an emap without key1, when add index by key1, it should fail.", func(emap EMap) {
			err := emap.Insert("key2", "value2")
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.AddIndex("key1", "index1")
			Expect(err).Should(HaveOccurred())
			Expect(emap.KeyNum()).To(BeEquivalentTo(1))
			Expect(emap.IndexNum()).To(BeEquivalentTo(0))
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an emap with key1 and index1, when add index1 to key1 again, it should fail.", func(emap EMap) {
			emap.Insert("key1", "value1", "index1")
			err := emap.AddIndex("key1", "index1")
			Expect(err).Should(HaveOccurred())
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an emap with key1, when add new index by key1, it should get key1's value by the new indices later.", func(emap EMap) {
			err := emap.Insert("key1", "value1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(emap.IndexNum()).To(BeEquivalentTo(0))

			err = emap.AddIndex("key1", "index1")
			Expect(err).ShouldNot(HaveOccurred())
			result1, err := emap.FetchByIndex("index1")
			Expect(result1).To(BeEquivalentTo([]interface{}{"value1"}))
			Expect(emap.KeyNum()).To(BeEquivalentTo(1))
			Expect(emap.IndexNum()).To(BeEquivalentTo(1))

			err = emap.AddIndex("key1", "index2")
			Expect(err).ShouldNot(HaveOccurred())
			result2, _ := emap.FetchByIndex("index2")
			Expect(result2).To(BeEquivalentTo([]interface{}{"value1"}))
			Expect(emap.KeyNum()).To(BeEquivalentTo(1))
			Expect(emap.IndexNum()).To(BeEquivalentTo(2))
			Expect(emap.IndexNumOfKey("key1")).To(BeEquivalentTo(2))
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an emap with key1 and index1, when remove index2 from key1, it should fail.", func(emap EMap) {
			emap.Insert("key1", "value1", "index1")
			err := emap.RemoveIndex("key1", "index2")
			Expect(err).Should(HaveOccurred())
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an emap with key1 and index1, when remove index from a non-existed key, it should fail.", func(emap EMap) {
			emap.Insert("key1", "value1", "index1")
			err := emap.RemoveIndex("key2", "index1")
			Expect(err).Should(HaveOccurred())
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)
	})

	Context("multi key and one index", func() {
		DescribeTable("Given an empty emap, when add values with different keys but same index, it should be able to get all values by the index.", func(emap EMap) {
			Expect(emap.HasIndex("index1")).To(Equal(false))
			err := emap.Insert("key1", "value1", "index1", "index2")
			Expect(emap.HasIndex("index1")).To(Equal(true))
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.Insert("key2", "value2", "index1")
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.Insert("key3", "value3", "index2")
			Expect(err).ShouldNot(HaveOccurred())

			result1, err := emap.FetchByIndex("index1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result1).To(BeEquivalentTo([]interface{}{"value1", "value2"}))
			result2, err := emap.FetchByIndex("index2")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result2).To(BeEquivalentTo([]interface{}{"value1", "value3"}))
			Expect(emap.KeyNum()).To(BeEquivalentTo(3))
			Expect(emap.IndexNum()).To(BeEquivalentTo(2))
			Expect(emap.KeyNumOfIndex("index1")).To(BeEquivalentTo(2))
			Expect(emap.KeyNumOfIndex("index2")).To(BeEquivalentTo(2))
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an emap with multi keys with same index, when delete index by key, it should not affect other keys.", func(emap EMap) {
			err := emap.Insert("key1", "value1", "index1", "index2")
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.Insert("key2", "value2", "index1")
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.Insert("key3", "value3", "index2")
			Expect(err).ShouldNot(HaveOccurred())

			result1, err := emap.FetchByIndex("index1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result1).To(BeEquivalentTo([]interface{}{"value1", "value2"}))

			err = emap.RemoveIndex("key1", "index1")
			Expect(err).ShouldNot(HaveOccurred())
			result2, err := emap.FetchByIndex("index1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result2).To(BeEquivalentTo([]interface{}{"value2"}))

			err = emap.RemoveIndex("key3", "index2")
			Expect(err).ShouldNot(HaveOccurred())
			result3, err := emap.FetchByIndex("index2")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result3).To(BeEquivalentTo([]interface{}{"value1"}))

			err = emap.RemoveIndex("key2", "index1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(emap.KeyNumOfIndex("index1")).To(BeEquivalentTo(0))
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an emap with multi keys and indices, when remove item by key, it should remove the related value and indices.", func(emap EMap) {
			err := emap.Insert("key1", "value1", "index1", "index2")
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.Insert("key2", "value2", "index2")
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.Insert("key3", "value3", "index2")
			Expect(err).ShouldNot(HaveOccurred())

			err = emap.DeleteByKey("key1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(emap.HasIndex("index1")).To(Equal(false))
			Expect(emap.HasIndex("index2")).To(Equal(true))
			_, err = emap.FetchByKey("key1")
			Expect(err).Should(HaveOccurred())
			_, err = emap.FetchByIndex("index1")
			Expect(err).Should(HaveOccurred())
			result1, err := emap.FetchByIndex("index2")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result1).To(BeEquivalentTo([]interface{}{"value2", "value3"}))
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an emap with multi keys and indices, when delete a item with a non-existed index, it should fail.", func(emap EMap) {
			err := emap.Insert("key1", "value1", "index1", "index2")
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.Insert("key2", "value2", "index2")
			Expect(err).ShouldNot(HaveOccurred())

			err = emap.DeleteByIndex("index3")
			Expect(err).Should(HaveOccurred())
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an emap with multi keys and indices, when remove item by index, it should remove all values related.", func(emap EMap) {
			err := emap.Insert("key1", "value1", "index1", "index2")
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.Insert("key2", "value2", "index2")
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.Insert("key3", "value3", "index1", "index3")
			Expect(err).ShouldNot(HaveOccurred())

			err = emap.DeleteByIndex("index2")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(emap.HasIndex("index2")).To(Equal(false))
			Expect(emap.HasKey("key1")).To(Equal(false))
			_, err = emap.FetchByIndex("index2")
			Expect(err).Should(HaveOccurred())
			_, err = emap.FetchByKey("key1")
			Expect(err).Should(HaveOccurred())
			result, err := emap.FetchByIndex("index1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result).To(BeEquivalentTo([]interface{}{"value3"}))
			Expect(emap.KeyNum()).To(BeEquivalentTo(1))
			Expect(emap.IndexNum()).To(BeEquivalentTo(2))
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an emap with multi keys and indices, when add a existed index to another value, it should be able to get all values by the index later.", func(emap EMap) {
			err := emap.Insert("key1", "value1", "index1", "index")
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.Insert("key2", "value2", "index2")
			Expect(err).ShouldNot(HaveOccurred())
			result, _ := emap.FetchByIndex("index")
			Expect(result).To(BeEquivalentTo([]interface{}{"value1"}))

			err = emap.AddIndex("key2", "index")
			Expect(err).ShouldNot(HaveOccurred())
			result, err = emap.FetchByIndex("index")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result).To(BeEquivalentTo([]interface{}{"value1", "value2"}))
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", "value", "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)
	})

	Context("expirable values", func() {
		var (
			emap EMap
		)

		BeforeEach(func() {
			emap = NewExpirableEMap(100)
		})

		AfterEach(func() {
			emap = nil
		})

		It("Given an expirable emap, when the value added has not IsExpired interface, it should fail.", func() {
			type testStruct struct {
				expired bool
			}

			Expect(emap.HasKey("key1")).To(Equal(false))
			value := new(testStruct)
			err := emap.Insert("key1", value, "index1", "index2", "index3")
			Expect(err).Should(HaveOccurred())
			Expect(emap.HasKey("key1")).To(Equal(false))
		})

		It("Given an interval to an expirable emap, when the value is expired, it should be collected.", func() {
			Expect(emap.HasKey("key1")).To(Equal(false))
			value := new(expirebleStruct)
			err := emap.Insert("key1", value, "index1", "index2", "index3")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(emap.HasKey("key1")).To(Equal(true))
			_, err = emap.FetchByKey("key1")
			Expect(err).ShouldNot(HaveOccurred())

			time.Sleep(time.Second)
			Expect(emap.HasKey("key1")).To(Equal(true))
			value.expired = true
			time.Sleep(time.Second)
			Expect(emap.HasKey("key1")).To(Equal(false))
		})
	})

	Context("strict emap", func() {
		type testStruct struct {
			data string
		}
		type anotherStruct struct {
			data string
		}

		BeforeEach(func() {
		})

		AfterEach(func() {
		})

		It("Given an unsupport type, when use this type to create a new strict emap, it should fail.", func() {
			_, err := NewStrictEMap([]string{"key"}, testStruct{"123"}, 123)
			Expect(err).Should(HaveOccurred())

			_, err = NewStrictEMap("key", testStruct{"123"}, map[string]string{})
			Expect(err).Should(HaveOccurred())
		})

		It("Given an empty strict emap, when use different types to add, it should fail.", func() {
			emap, err := NewStrictEMap("key", testStruct{"123"}, 123)
			Expect(err).ShouldNot(HaveOccurred())

			err = emap.Insert("key", testStruct{"123"}, "123")
			Expect(err).Should(HaveOccurred())

			err = emap.Insert(123, testStruct{"123"}, 123)
			Expect(err).Should(HaveOccurred())

			err = emap.Insert("sample", "test", 123)
			Expect(err).Should(HaveOccurred())

			err = emap.Insert("sample", anotherStruct{"123"}, 123)
			Expect(err).Should(HaveOccurred())

			Expect(emap.KeyNum()).To(BeEquivalentTo(0))
			Expect(emap.IndexNum()).To(BeEquivalentTo(0))
		})

		It("Given an strict emap, when use its interfaces, it should check the type of input paras.", func() {
			emap, err := NewStrictEMap("key", testStruct{"123"}, 123)
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.Insert("key", testStruct{"123"}, 123)
			Expect(err).ShouldNot(HaveOccurred())

			err = emap.DeleteByKey(123)
			Expect(err).Should(HaveOccurred())
			err = emap.DeleteByIndex("123")
			Expect(err).Should(HaveOccurred())
			_, err = emap.FetchByKey(123)
			Expect(err).Should(HaveOccurred())
			_, err = emap.FetchByIndex("123")
			Expect(err).Should(HaveOccurred())
			err = emap.AddIndex("key", "123")
			Expect(err).Should(HaveOccurred())
			err = emap.AddIndex(123, 123)
			Expect(err).Should(HaveOccurred())
			err = emap.RemoveIndex("key", "123")
			Expect(err).Should(HaveOccurred())
			err = emap.RemoveIndex(123, 123)
			Expect(err).Should(HaveOccurred())
			Expect(emap.IndexNumOfKey(123)).Should(Equal(0))
			Expect(emap.KeyNumOfIndex("123")).Should(Equal(0))
		})

		It("Given an strict emap, when use HasKey or HasIndex interface with different key or index type, it should return false.", func() {
			emap, err := NewStrictEMap("key", testStruct{"123"}, 123)
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.Insert("key", testStruct{"123"}, 123)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(emap.HasKey(123)).Should(Equal(false))
			Expect(emap.HasIndex("123")).Should(Equal(false))
		})
	})

	Context("Higher-order functions", func() {
		type testStruct struct {
			data string
		}

		type anotherStruct struct {
			data string
		}

		DescribeTable("Given an emap, when call Transform interface, it should return the trasformed values related to the callback.", func(emap EMap) {
			emap.Insert("key1", 1, "index1")
			emap.Insert("key2", 2, "index2")
			emap.Insert("key3", 3, "index3")
			Expect(emap.KeyNum()).Should(BeEquivalentTo(3))

			callback := func(key interface{}, value interface{}) (interface{}, error) {
				return value.(int) + 10, nil
			}
			targets, err := emap.Transform(callback)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(targets)).Should(BeEquivalentTo(emap.KeyNum()))
			Expect(targets["key1"]).Should(BeEquivalentTo(11))
			Expect(targets["key2"]).Should(BeEquivalentTo(12))
			Expect(targets["key3"]).Should(BeEquivalentTo(13))
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", 0, "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an emap, when call Transform interface, it should fail when callback fails.", func(emap EMap) {
			emap.Insert("key1", 1, "index1")
			emap.Insert("key2", 2, "index2")
			emap.Insert("key3", 3, "index3")
			emap.Insert("keyError", 123)
			Expect(emap.KeyNum()).Should(BeEquivalentTo(4))

			callback := func(key interface{}, value interface{}) (interface{}, error) {
				if key == "keyError" {
					return nil, errors.New("error key")
				}
				return value.(int) + 10, nil
			}
			targets, err := emap.Transform(callback)
			Expect(err).Should(HaveOccurred())
			Expect(len(targets)).Should(BeEquivalentTo(0))
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", 0, "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)

		DescribeTable("Given an emap, when call Foreach interface, it should apply callback to each item.", func(emap EMap) {
			type testStruct struct {
				num int
			}
			emap.Insert("key1", &testStruct{1}, "index1")
			emap.Insert("key2", &testStruct{2}, "index2")
			emap.Insert("key3", &testStruct{3}, "index3")
			Expect(emap.KeyNum()).Should(BeEquivalentTo(3))

			callback := func(key interface{}, value interface{}) {
				value.(*testStruct).num = value.(*testStruct).num + 10
			}
			emap.Foreach(callback)
			Expect(emap.FetchByKey("key1")).To(BeEquivalentTo(&testStruct{11}))
			Expect(emap.FetchByKey("key2")).To(BeEquivalentTo(&testStruct{12}))
			Expect(emap.FetchByKey("key3")).To(BeEquivalentTo(&testStruct{13}))
		},
			Entry("generic emap test", NewGenericEMap()),
			Entry("strict emap test", NewStrictEmapWrapper("key", &testStruct{"value"}, "index")),
			Entry("nolock emap test", NewUnlockEMap()),
		)
	})

	Context("benchmark emap", func() {
		BeforeEach(func() {
		})

		AfterEach(func() {
		})

		Measure("Benchmark the goalng map performance", func(b Benchmarker) {
			goMap := make(map[interface{}]interface{})
			b.Time("GoMap", func() {
				GoMapAdd(goMap, 200000)
				GoMapGet(goMap, 200000)
				GoMapDel(goMap, 200000)
			})
		}, 10)

		Measure("Benchmark the generic emap performance", func(b Benchmarker) {
			eMap := NewGenericEMap()
			EMapRuntime := b.Time("GenericEMap", func() {
				EMapAdd(eMap, 200000)
				EMapGet(eMap, 200000)
				EMapDel(eMap, 200000)
			})

			Ω(EMapRuntime.Seconds()).Should(BeNumerically("<", 2), "Add/Get/Del 200000 values shouldn't take too long.")
		}, 10)

		Measure("Benchmark the stric emap performance", func(b Benchmarker) {
			eMap, _ := NewStrictEMap("sample", &expirebleStruct{false, 0}, 123)
			EMapRuntime := b.Time("StrictEMap", func() {
				EMapAdd(eMap, 200000)
				EMapGet(eMap, 200000)
				EMapDel(eMap, 200000)
			})

			Ω(EMapRuntime.Seconds()).Should(BeNumerically("<", 2), "Add/Get/Del 200000 values shouldn't take too long.")
		}, 10)

		Measure("Benchmark the nolock emap performance", func(b Benchmarker) {
			eMap := NewUnlockEMap()
			EMapRuntime := b.Time("NolockEMap", func() {
				EMapAdd(eMap, 200000)
				EMapGet(eMap, 200000)
				EMapDel(eMap, 200000)
			})

			Ω(EMapRuntime.Seconds()).Should(BeNumerically("<", 2), "Add/Get/Del 200000 values shouldn't take too long.")
		}, 10)
	})

	Context("debug", func() {
		var (
			emap *GenericEMap
		)

		BeforeEach(func() {
			emap = NewGenericEMap()
		})

		AfterEach(func() {
		})

		It("Fix the unconsistency caused by deleteByKey", func() {
			emap.Insert("key1", "value1", "index1", "index2", "index3")
			emap.Insert("key2", "value2", "index1", "index2", "index3")
			emap.Insert("key3", "value3", "index1", "index2", "index3")
			Expect(emap.check()).ShouldNot(HaveOccurred())

			emap.DeleteByKey("key1")
			Expect(emap.check()).ShouldNot(HaveOccurred())
		})

		It("Fix the unconsistency caused by deleteByIndex", func() {
			emap.Insert("key1", "value1", "index1", "index2", "index3")
			emap.Insert("key2", "value2", "index1", "index2", "index3")
			emap.Insert("key3", "value3", "index1", "index2", "index3")
			Expect(emap.check()).ShouldNot(HaveOccurred())

			emap.DeleteByIndex("index1")
			Expect(emap.check()).ShouldNot(HaveOccurred())
		})
	})
})

func NewStrictEmapWrapper(key interface{}, value interface{}, index interface{}) (emap EMap) {
	emap, _ = NewStrictEMap(key, value, index)
	return
}

type expirebleStruct struct {
	expired bool
	number  int
}

func (v *expirebleStruct) IsExpired() bool {
	return v.expired
}

func GoMapAdd(goMap map[interface{}]interface{}, number int) {
	for i := 0; i < number; i++ {
		goMap[string(i)] = &expirebleStruct{false, i}
	}
}

func GoMapGet(goMap map[interface{}]interface{}, number int) (dump interface{}) {
	for i := 0; i < number; i++ {
		dump = goMap[string(i)].(*expirebleStruct).number
	}
	return
}

func GoMapDel(goMap map[interface{}]interface{}, number int) {
	for i := 0; i < number; i++ {
		delete(goMap, string(i))
	}
}

func EMapAdd(emap EMap, number int) {
	for i := 0; i < number; i++ {
		emap.Insert(string(i), &expirebleStruct{false, i})
	}
}

func EMapGet(emap EMap, number int) (dump interface{}) {
	for i := 0; i < number; i++ {
		value, _ := emap.FetchByKey(string(i))
		dump = value.(*expirebleStruct).number
	}
	return
}

func EMapDel(emap EMap, number int) {
	for i := 0; i < number; i++ {
		emap.DeleteByKey(string(i))
	}
}
