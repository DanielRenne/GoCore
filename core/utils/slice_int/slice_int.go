package slice_int

func Init(elements ...int) []int {
	var int_slice = Allocate()
	for _, element := range elements {
		int_slice = Extend(int_slice, element)
	}
	return int_slice
}

func Allocate() []int {
	return make([]int, 0, 1)
}

func Extend(slice []int, elements ...int) []int {
	for _, element := range elements {
		n := len(slice)
		if n == cap(slice) {
			// Slice is full; must grow.
			// We double its size and add 1, so if the size is zero we still grow.
			newSlice := make([]int, len(slice), 2*len(slice)+1)
			copy(newSlice, slice)
			slice = newSlice
		}
		slice = slice[0 : n+1]
		slice[n] = element
	}
	return slice
}
