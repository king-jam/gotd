package giphy

import (
	"fmt"
	"log"
	"testing"
)

func TestTags(t *testing.T) {
	u := "https://giphy.com/gifs/stan-wawrinka-120sn8DzgO12Yo/"
	tags, err := GetGIFTags(u)
	if err != nil {
		log.Fatal("AAHHHHH 2")
	}
	for _, t := range tags {
		fmt.Printf("%s\n", t)
	}
	fmt.Print(tags)
}
