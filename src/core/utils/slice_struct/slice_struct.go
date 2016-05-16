package slice_struct

func Init() []interface{} {
    return make([]interface{}, 0, 1)
}

func Extend(slice []interface{}, element interface{}) []interface{} {
    n := len(slice)
    if n == cap(slice) {
        // Slice is full; must grow.
        // We double its size and add 1, so if the size is zero we still grow.
        newSlice := make([]interface{}, len(slice), 2*len(slice)+1)
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[0 : n+1]
    slice[n] = element
    return slice
}