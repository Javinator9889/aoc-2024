package disk

import "github.com/Javinator9889/aoc-2024/2024/day09/block"

type Disk []*block.Block

func (d Disk) FreeBlock() (int, *block.Block) {
	for i, b := range d {
		if b.Free() > 0 {
			return i, b
		}
	}
	return -1, nil
}
