// Package zuint32 implements radix sort for []uint32.
package zuint32

import (
	"sort"
)

const (
	// MinSize is the minimum size of a slice that will be radix sorted by Sort.
	MinSize      = 128
	radix   uint = 8
	bitSize uint = 32
)

// Sort sorts x using a Radix sort (Small slices are sorted with sort.Sort() instead).
func Sort(x []uint32) {
	if len(x) < MinSize {
		sort.Slice(x, func(i, j int) bool { return x[i] < x[j] })
	} else {
		buffer := make([]uint32, len(x))
		SortBYOB(x, buffer)
	}
}

// SortCopy is similar to Sort, but returns a sorted copy of x, leaving x unmodified.
func SortCopy(x []uint32) []uint32 {
	y := make([]uint32, len(x))
	copy(y, x)
	Sort(y)
	return y
}

// SortBYOB sorts x using a Radix sort, using supplied buffer space. Panics if
// len(x) does not equal len(buffer). Uses radix sort even on small slices.
func SortBYOB(x, buffer []uint32) {
	if len(x) > len(buffer) {
		panic("Buffer too small")
	}
	if len(x) < 2 {
		return
	}

	from := x
	to := buffer[:len(x)]
	var key uint8       // Current byte value
	var offset [256]int // Keep track of where room is made for byte groups in the buffer

	for keyOffset := uint(0); keyOffset < bitSize; keyOffset += radix {
		keyMask := uint32(0xFF << keyOffset)
		var counts [256]int // Keep track of the number of elements for each kind of byte
		sorted := false
		prev := uint32(0)

		for _, elem := range from {
			// For each elem to sort, fetch the byte at current radix
			key = uint8((elem & keyMask) >> keyOffset)
			// inc count of bytes of this type
			counts[key]++
			if sorted { // Detect sorted
				sorted = elem >= prev
				prev = elem
			}
		}

		if sorted { // Short-circuit sorted
			if (keyOffset/radix)%2 == 1 {
				copy(to, from)
			}
			return
		}

		// Make room for each group of bytes by calculating offset of each
		offset[0] = 0
		for i := 1; i < len(offset); i++ {
			offset[i] = offset[i-1] + counts[i-1]
		}

		// Swap values between the buffers by radix
		for _, elem := range from {
			key = uint8((elem & keyMask) >> keyOffset) // Get the byte of each element at the radix
			to[offset[key]] = elem                     // Copy the element depending on byte offsets
			offset[key]++                              // One less space, move the offset
		}
		// Each pass copy data the other way
		to, from = from, to
	}
}
