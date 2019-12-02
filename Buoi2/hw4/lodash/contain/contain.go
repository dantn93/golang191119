package contain

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/patrickmn/go-cache"
)

func Contain(c *cache.Cache, intSlice *[]int, element int, t time.Duration) interface{} {
	byteListNumber := []byte{}
	l, found := c.Get("intSlice")
	if found {
		byteList, err := json.Marshal(l)
		if err != nil {
			log.Println("Can not parse to byte")
			panic(err)
		}
		byteListNumber = byteList
	}

	byteIntSlice, err := json.Marshal(intSlice)
	if err != nil {
		log.Println("Can not parse to byte")
		panic(err)
	}

	//Check slice in caching is the same as input
	if cmp.Equal(byteListNumber, byteIntSlice) {
		e, found := c.Get("element")
		if found && e == element {
			result, found := c.Get("result")
			if found {
				return result
			}
		}
		result := FindElementInSlice(intSlice, element)
		Caching(c, intSlice, element, result, t)
		return result
	} else {
		result := FindElementInSlice(intSlice, element)
		Caching(c, intSlice, element, result, t)
		return result
	}
}

func FindElementInSlice(intSlice *[]int, element int) bool {
	for _, v := range *intSlice {
		if v == element {
			return true
		}
	}
	return false
}

func Caching(c *cache.Cache, intSlice *[]int, element int, result bool, t time.Duration) {
	c.Set("intSlice", *intSlice, t)
	c.Set("element", element, t)
	c.Set("result", result, t)
}
