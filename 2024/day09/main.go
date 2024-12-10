package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Javinator9889/aoc-2024/2024/day09/block"
	"github.com/Javinator9889/aoc-2024/cast"
	"github.com/Javinator9889/aoc-2024/util"
)

//go:embed input.txt
var input string

func init() {
	// do this in init (not main) so test file has same input
	input = strings.TrimRight(input, "\n")
	if len(input) == 0 {
		panic("empty input.txt file")
	}
}

func main() {
	var part int
	var debug bool
	flag.IntVar(&part, "part", 1, "part 1 or 2")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()
	fmt.Println("Running part", part)
	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if part == 1 {
		ans := part1(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

type Disk []*block.Block

func (d Disk) FreeBlock(minSize int) (int, *block.Block) {
	for i, b := range d {
		if b.Free() >= minSize {
			return i, b
		}
	}
	return -1, nil
}

func (d Disk) String() string {
	str := strings.Builder{}
	for _, b := range d {
		str.WriteString(b.String())
	}
	return str.String()
}

func part1(input string) (chksum int) {
	disk := parseInput(input)
	slog.Debug("Disk layout", "disk", disk)
	// Reallocate the blocks, starting from the end. As we're moving the files from the end,
	// we can just reallocate until the middle of the disk.
outer:
	for i := len(disk) - 1; i >= 0; i-- {
		blk := disk[i]
		if blk.Free() == blk.MaxSize {
			// Skip empty blocks
			continue
		}
		// Reallocate the block
		freeIdx, freeBlk := disk.FreeBlock(1)
		slog.Debug("Reallocating", "block", blk)
		for f := blk.PopLast(); f != nil; {
			if freeIdx == -1 {
				blk.Add(f)
				slog.Debug("No free blocks left")
				break outer
			}
			// Do not move a block after its position
			if freeIdx >= i {
				blk.Add(f)
				slog.Debug("No free blocks after current block", "block", i)
				break outer
			}
			if freeBlk.AddPartial(f) {
				break
			}
			freeIdx, freeBlk = disk.FreeBlock(1)
		}
		slog.Debug("Disk after reallocating", "i", i, "disk", disk)
	}
	slog.Debug("Final disk", "disk", disk)

	// Compute the new checksum
	offset := 0
	for _, blk := range disk {
		chksum += blk.Chksum(&offset)
	}
	return
}

func part2(input string) (chksum int) {
	// Part two is the same as part one, but we need to find enough room for a whole file to fit
	// on a free block.
	disk := parseInput(input)
	slog.Debug("Disk layout", "disk", disk)
	for i := len(disk) - 1; i >= 0; i-- {
		blk := disk[i]
		if blk.Free() == blk.MaxSize {
			// Skip empty blocks
			continue
		}
		freeIdx, freeBlk := disk.FreeBlock(blk.Length())
		if freeIdx == -1 || freeIdx >= i {
			// No free blocks left or no blocks after the current one
			slog.Debug("No room for re-arranging block", "block", blk)
			continue
		}
		// Reallocate the block
		slog.Debug("Reallocating", "block", blk)
		for {
			var f *block.File
			if f = blk.PopLast(); f == nil {
				break
			}
			// Add the file to the free block
			if !freeBlk.Add(f) {
				// The file should fit, there's an error
				slog.Error("File does not fit in free block", "file", f, "block", freeBlk)
				panic("File does not fit in free block")
			}
			// Move to the next free block
			freeIdx, freeBlk = disk.FreeBlock(blk.Length())
			if freeIdx == -1 || freeIdx >= i {
				// No free blocks left or no blocks after the current one
				break
			}
		}
		slog.Debug("Disk after reallocating", "i", i, "disk", disk)
	}
	slog.Debug("Final disk", "disk", disk)

	// Compute the new checksum
	offset := 0
	for _, blk := range disk {
		chksum += blk.Chksum(&offset)
	}
	return
}

func parseInput(input string) Disk {
	disk := make(Disk, 0)
	idx := 0
	// The input is a single line with a series of numbers
	for i, c := range input {
		if c == '\n' || c == '0' { // Skip newlines and zeroes
			continue
		}

		slog.Debug("Parsing input", "i", i, "c", c, "char", string(c))
		blk := &block.Block{}
		length := cast.ToInt(string(c))
		if i%2 == 0 {
			file := &block.File{ID: idx, Length: length}
			blk.Push(file)
			idx++
		} else {
			blk.MaxSize = length
		}
		disk = append(disk, blk)
	}
	return disk
}
