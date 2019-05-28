package gif

import (
	"fmt"
	"testing"
)

func TestTags(t *testing.T) {
	gif := GIF{
		GIF: "monday morning",
	}
	err := BuildGif(&gif)
	if err != nil {
		fmt.Println("FREAK OUT")
	}

}
