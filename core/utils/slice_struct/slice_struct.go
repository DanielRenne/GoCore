package slice_struct

func Init(elements ...interface{}) []interface{} {
    var struct_slice = Allocate()
    for _, element := range elements {
        struct_slice = Extend(struct_slice, element)
    }
    return struct_slice
}

func Allocate() []interface{} {
    return make([]interface{}, 0, 1)
}

func Extend(slice []interface{}, elements ...interface{}) []interface{} {
    for _, element := range elements {
        n := len(slice)
        if n == cap(slice) {
            // Slice is full; must grow.
            // We double its size and add 1, so if the size is zero we still grow.
            newSlice := make([]interface{}, len(slice), 2 * len(slice) + 1)
            copy(newSlice, slice)
            slice = newSlice
        }
        slice = slice[0 : n + 1]
        slice[n] = element
    }
    return slice
}