package block

import (
	"strings"

	"github.com/Javinator9889/aoc-2024/cast"
)

type File struct {
	ID     int
	Length int
}

func (f File) String() string {
	str := strings.Builder{}
	for i := 0; i < f.Length; i++ {
		str.WriteString(cast.ToString(f.ID))
	}
	return str.String()
}

type Block struct {
	MaxSize int
	files   []*File
}

func (b Block) Length() (length int) {
	for _, f := range b.files {
		if f == nil {
			continue
		}
		length += f.Length
	}
	return
}

func (b Block) Free() int {
	return b.MaxSize - b.Length()
}

// Adds a file, increasing the block's size
func (b *Block) Push(file *File) int {
	b.files = append(b.files, file)
	b.MaxSize += file.Length
	return file.Length
}

func (b *Block) Add(file *File) bool {
	if b.Free() < file.Length {
		return false
	}
	b.files = append(b.files, file)
	return true
}

func (b *Block) AddPartial(file *File) bool {
	free := b.Free()
	if free == 0 {
		return false
	}
	needle := min(free, file.Length)
	// Fast path: The file fits in the block
	if needle == file.Length {
		b.files = append(b.files, file)
		return true
	}
	partial := &File{ID: file.ID, Length: needle}
	b.files = append(b.files, partial)
	file.Length -= needle
	return false
}

func (b *Block) PopLast() *File {
	if len(b.files) == 0 {
		return nil
	}
	file := b.files[len(b.files)-1]
	b.files = b.files[:len(b.files)-1]
	return file
}

func (b Block) Chksum(offset *int) (chksum int) {
	base := *offset
	for _, f := range b.files {
		if f == nil {
			continue
		}
		for i := 0; i < f.Length; i++ {
			chksum += f.ID * base
			base++
		}
	}
	// The offset is increased by the block's size
	*offset += b.MaxSize
	return
}

func (b Block) String() string {
	str := strings.Builder{}
	for _, f := range b.files {
		str.WriteString(f.String())
	}
	for i := 0; i < b.Free(); i++ {
		str.WriteRune('.')
	}
	return str.String()
}
