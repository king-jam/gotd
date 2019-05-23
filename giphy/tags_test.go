package giphy

import (
	"fmt"
	"testing"
)

func TestTags(t *testing.T) {
	tags, err := GetGIFTags("https://giphy.com/gifs/prepareforwinter-l4FGkZk8tCWDVPsB2/fullscreen")
	if err != nil {
		fmt.Println("FREAK OUT")
	}
	for _, t := range tags {
		fmt.Printf("%s\n", t)
	}

}
