package emap

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Tests of emap", func() {
	var (
		emap EMap
	)

	Context("one unique key and multi indices", func() {
		BeforeEach(func() {
			emap = NewGenericEMap()
		})

		AfterEach(func() {
			emap = nil
		})

		It("Given an empty emap, when add a new item, it should be able to get by key or index later.", func() {
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
		})

		It("Given an emap with key1, when add a new item with the same key1, it should fail.", func() {
			err := emap.Insert("key1", "value1", "index1", "index2")
			Expect(err).ShouldNot(HaveOccurred())

			err = emap.Insert("key1", "value2", "index3")
			Expect(err).Should(HaveOccurred())

			Expect(emap.KeyNum()).To(BeEquivalentTo(1))
			Expect(emap.IndexNum()).To(BeEquivalentTo(2))
			Expect(emap.IndexNumOfKey("key1")).To(BeEquivalentTo(2))
		})

		It("Given an empty emap, when delete a item with the key1, it should fail.", func() {
			err := emap.DeleteByKey("key1")
			Expect(err).Should(HaveOccurred())
		})

		It("Given an emap with multi values, when delete by key, it should delete the value and indices of the key.", func() {
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
		})

		It("Given an emap without key1, when add index by key1, it should fail.", func() {
			err := emap.Insert("key2", "value2")
			Expect(err).ShouldNot(HaveOccurred())
			err = emap.AddIndex("key1", "index1")
			Expect(err).Should(HaveOccurred())
			Expect(emap.KeyNum()).To(BeEquivalentTo(1))
			Expect(emap.IndexNum()).To(BeEquivalentTo(0))
		})

		It("Given an emap with key1, when add new indices by key1, it should get key1's value by the new indices later.", func() {
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
			result2, err := emap.FetchByIndex("index2")
			Expect(result2).To(BeEquivalentTo([]interface{}{"value1"}))
			Expect(emap.KeyNum()).To(BeEquivalentTo(1))
			Expect(emap.IndexNum()).To(BeEquivalentTo(2))
			Expect(emap.IndexNumOfKey("key1")).To(BeEquivalentTo(2))
		})

	})

	Context("multi key and one indices", func() {
		BeforeEach(func() {
			emap = NewGenericEMap()
		})

		AfterEach(func() {
			emap = nil
		})

		It("Given an empty emap, when add values with different keys but same index, it should be able to get all values by the index.", func() {
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
		})

		It("Given an emap with multi keys with same index, when delete index by key, it should not affect other keys.", func() {
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
		})

		It("Given an emap with multi keys and indices, when remove item by key, it should remove the related value and indices.", func() {
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
		})

		It("Given an emap with multi keys and indices, when remove item by index, it should remove all values related.", func() {
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
			_, err = emap.FetchByKey(123)
			Expect(err).Should(HaveOccurred())
			_, err = emap.FetchByIndex("123")
			Expect(err).Should(HaveOccurred())
			err = emap.AddIndex("key", "123")
			Expect(err).Should(HaveOccurred())
			err = emap.RemoveIndex(123, 123)
			Expect(err).Should(HaveOccurred())
			Expect(emap.IndexNumOfKey(123)).Should(Equal(0))
			Expect(emap.KeyNumOfIndex("123")).Should(Equal(0))
		})
	})

	Context("expirable emap", func() {
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

	Context("nolock emap", func() {
		BeforeEach(func() {
			emap = NewUnlockEMap()
		})

		AfterEach(func() {
			emap = nil
		})

		It("Given an empty emap, when add values with different keys but same index, it should be able to get all values by the index.", func() {
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
		})

		It("Given an emap with multi keys with same index, when delete index by key, it should not affect other keys.", func() {
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
		})

		It("Given an emap with multi keys and indices, when remove item by key, it should remove the key from its indices.", func() {
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
			_, err = emap.FetchByIndex("index1")
			Expect(err).Should(HaveOccurred())
			result1, err := emap.FetchByIndex("index2")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result1).To(BeEquivalentTo([]interface{}{"value2", "value3"}))
		})
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

		Measure("Benchmark the expirable emap performance", func(b Benchmarker) {
			eMap := NewExpirableEMap(1000)
			EMapRuntime := b.Time("ExpirableEMap", func() {
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
})

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
