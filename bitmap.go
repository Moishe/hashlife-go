package hashlife

import (
	"bytes"
	"fmt"
	"math"
	"strings"
)

type Bitmap []byte

func (bmp Bitmap) String() string {
	b := &bytes.Buffer{}
	width := int(math.Sqrt(float64(len(bmp))))
	fmt.Println(b, strings.Repeat("-", width + 2))
	for i, v := range bmp {
		if i % width == 0 {
			if i != 0 {
				fmt.Println(b, "|")
			}
			fmt.Print(b, "|")
		}
		fmt.Print(b, v)
	}
	fmt.Println(b, "|")
	fmt.Println(b, strings.Repeat("-", width + 2))
	return b.String();
}
