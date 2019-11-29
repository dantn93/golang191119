package contain

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/patrickmn/go-cache"
)

func Contain(c *cache.Cache, intSlice *[]int, element int, t time.Duration) bool {
	defer Caching(c, intSlice, element, t)

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
		iresult, found := c.Get("element")
		if found && iresult == element {
			return true
		}
		return FindElementInSlice(intSlice, element)
	} else {
		return FindElementInSlice(intSlice, element)
	}
}

func FindElementInSlice(intSlice *[]int, element int) bool {
	for _, v := range *intSlice {
		time.Sleep(10000)
		if v == element {
			return true
		}
	}
	return false
}

func Caching(c *cache.Cache, intSlice *[]int, element int, t time.Duration) {
	c.Set("intSlice", *intSlice, t)
	c.Set("element", element, t)
}
