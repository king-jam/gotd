package gif

import (
	"fmt"
	"testing"
)

func TestTags(t *testing.T) {
	gif := GIF{
		GIF: "lord of the ring",
	}
	err := BuildGif(&gif)
	if err != nil {
		fmt.Println("FREAK OUT")
	}

}
