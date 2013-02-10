package hashlife

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"testing"
)

var Benchmark_Level = flag.Int("benchmark_level", 10, "The exponent (2**e) of the number of cells per side of the benchmark board.")

func AssertEqual(x byte, y byte, t *testing.T) {
	if x != y {
		t.Error("Mismatch")
	}
}

func TestTreeFromBit(t *testing.T) {
	bitmap := []byte{1}
	leaf, err := TreeFromBitmapBase(bitmap)

	if err != nil {
		t.Error("TreeFromBitmapBase returned an error: %s", err)
	}

	if leaf == nil {
		t.Error("TreeFromBitmap with 1 cell should return leaf.")
	}
}

func TestCountsFromLeaves(t *testing.T) {
	bitmap := []byte{
		1, 1, 1, 0,
		1, 0, 1, 0,
		1, 1, 1, 0,
		0, 0, 0, 0,
	}
	root, _ := TreeFromBitmapBase(bitmap)
	ulc, urc, llc, lrc := CountsFromLeaves(root)
	AssertEqual(8, ulc, t)
	AssertEqual(4, urc, t)
	AssertEqual(4, llc, t)
	AssertEqual(2, lrc, t)
	// TODO(moishel): add more comprehensive testing here.
}

func TestTreeFromBitmapBase(t *testing.T) {
	width := 1 << 5
	bitmap := make([]byte, width*width)
	for i := 0; i < width*width; i++ {
		if rand.Float32() < 0.5 {
			bitmap[i] = 1
		}
	}
	root, _ := TreeFromBitmapBase(bitmap)
	dumped := DumpNode(root)
	if !bytes.Equal(bitmap, dumped) {
		t.Error("Dumped bitmaps don't match.")
	}
}

func MakeNBitmap(n int) []byte {
	bitmap := make([]byte, n*n)
	for i := 0; i < n*n; i++ {
		if rand.Float32() < 0.5 {
			bitmap[i] = 1
		}
	}

	return bitmap
}

func TestTreeFromNGrid(t *testing.T) {
	for j := 3; j < 5; j++ {
		bitmap := MakeNBitmap(1 << uint(j))
		root, _ := TreeFromBitmapBase(bitmap)
		newbitmap := DumpNode(root)

		for i := 0; i < len(bitmap); i++ {
			if bitmap[i] != newbitmap[i] {
				t.Error("Mismatched bitmaps")
			}
		}
	}
}

func TestNextGenerationLevel3(t *testing.T) {
	bitmap := []byte{
		0, 0, 0, 0,
		0, 1, 1, 0,
		1, 1, 0, 0,
		0, 1, 0, 0,
	}

	root, _ := TreeFromBitmapBase(bitmap)
	result := NextGeneration(root)
	simple_result := make([]byte, 4)
	SimpleNextGeneration(&bitmap, &simple_result, 4)
	dump_result := DumpNode(result)
	if len(simple_result) != len(dump_result) {
		t.Error("Lengths of results from board-size 4x4 don't match")
	}
	for i := 0; i < len(dump_result); i++ {
		if simple_result[i] != dump_result[i] {
			t.Error("Invalid result")
		}
	}
}

func TestNextGenerationLevel4(t *testing.T) {
	bitmap := []byte{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	VerifyIdenticalResults(bitmap, t)
}

func VerifyIdenticalResults(bitmap []byte, t *testing.T) {
	root, _ := TreeFromBitmapBase(bitmap)
	next := NextGeneration(root)
	hash_dump := DumpNode(next)
	simple_next := SimpleNthGeneration(bitmap)

	if len(hash_dump) != len(simple_next) {
		t.Error("Length of dumped bitmaps don't match.")
	}

	if !bytes.Equal(hash_dump, simple_next) {
		t.Error("Dumped bitmaps don't match.")
		fmt.Println(root)
		DumpBitmap(DumpNode(root))
		DumpBitmap(bitmap)
		DumpBitmap(hash_dump)
		DumpBitmap(simple_next)
	}
}

func TestCenteredRPentomino(t *testing.T) {
	bitmap := CenteredRPentomino(3)
	VerifyIdenticalResults(bitmap, t)
}

func TestCenteredRPentominoLarge(t *testing.T) {
	bitmap := CenteredRPentomino(10)
	VerifyIdenticalResults(bitmap, t)
}

func TestRandomBitmap(t *testing.T) {
	for level := 9; level < 11; level++ {
		width := 1 << uint(level)
		bitmap := make([]byte, width*width)
		for i := 0; i < width*width; i++ {
			if rand.Float32() < 0.5 {
				bitmap[i] = 1
			}
		}
		VerifyIdenticalResults(bitmap, t)
	}
}

func CenteredRPentomino(level int) []byte {
	width := 1 << uint(level)
	size := width * width
	bitmap := make([]byte, size)

	middle := width/2 + (size / 2)
	bitmap[middle] = 1
	bitmap[middle-1] = 1
	bitmap[middle+1] = 1
	bitmap[middle+width] = 1
	bitmap[middle-width+1] = 1

	return bitmap
}

func RandomBitmap(level int) []byte {
	width := 1 << uint(level)
	bitmap := make([]byte, width*width)
	for i := 0; i < width*width; i++ {
		if rand.Float32() < 0.2 {
			bitmap[i] = 1
		}
	}
	return bitmap
}

func RunBenchmarkWithBitmap(bitmap []byte, b *testing.B) {
	var root *Node
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		nodeMap = make(NodeMap)
		root, _ = TreeFromBitmapBase(bitmap)
		b.StartTimer()
		NextGeneration(root)
	}
}

func Benchmark_RPentominoNextGeneration(b *testing.B) {
	b.StopTimer()
	bitmap := CenteredRPentomino(*Benchmark_Level)
	RunBenchmarkWithBitmap(bitmap, b)
}

func Benchmark_RandomNextGeneration(b *testing.B) {
	b.StopTimer()
	bitmap := RandomBitmap(*Benchmark_Level)
	RunBenchmarkWithBitmap(bitmap, b)
}

func Benchmark_RPentominoSimpleNextGeneration(b *testing.B) {
	b.StopTimer()
	bitmap := CenteredRPentomino(*Benchmark_Level)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		SimpleNthGeneration(bitmap)
	}
}

func Benchmark_RandomSimpleNextGeneration(b *testing.B) {
	b.StopTimer()
	bitmap := RandomBitmap(*Benchmark_Level)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		SimpleNthGeneration(bitmap)
	}
}
