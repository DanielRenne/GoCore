package slice_string

func Init(elements ...string) []string {
    var str_slice = Allocate()
    for _, element := range elements {
        str_slice = Extend(str_slice, element)
    }
    return str_slice
}

func Allocate() []string {
    return make([]string, 0, 1)
}

func Extend(slice []string, elements ...string) []string {
    for _, element := range elements {
        n := len(slice)
        if n == cap(slice) {
            // Slice is full; must grow.
            // We double its size and add 1, so if the size is zero we still grow.
            newSlice := make([]string, len(slice), 2*len(slice)+1)
            copy(newSlice, slice)
            slice = newSlice
        }
        slice = slice[0 : n+1]
        slice[n] = element
    }
    return slice
}