package giphy

import (
	"fmt"
	"log"
	"testing"
)

func TestTags(t *testing.T) {
	u := "https://giphy.com/gifs/prepareforwinter-l4FGkZk8tCWDVPsB2/fullscreen"
	tags, err := GetGIFTags(u)
	if err != nil {
		log.Fatal("AAHHHHH 2")
	}
	for _, t := range tags {
		fmt.Printf("%s\n", t)
	}

}
