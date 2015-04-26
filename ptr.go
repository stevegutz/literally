package literally

// BoolPtr returns a pointer to the provided value
func BoolPtr(v bool) *bool { return &v }

// IntPtr returns a pointer to the provided value
func IntPtr(v int) *int { return &v }

// Int8Ptr returns a pointer to the provided value
func Int8Ptr(v int8) *int8 { return &v }

// Int16Ptr returns a pointer to the provided value
func Int16Ptr(v int16) *int16 { return &v }

// Int32Ptr returns a pointer to the provided value
func Int32Ptr(v int32) *int32 { return &v }

// Int64Ptr returns a pointer to the provided value
func Int64Ptr(v int64) *int64 { return &v }

// UintPtr returns a pointer to the provided value
func UintPtr(v uint) *uint { return &v }

// Uint8Ptr returns a pointer to the provided value
func Uint8Ptr(v uint8) *uint8 { return &v }

// Uint16Ptr returns a pointer to the provided value
func Uint16Ptr(v uint16) *uint16 { return &v }

// Uint32Ptr returns a pointer to the provided value
func Uint32Ptr(v uint32) *uint32 { return &v }

// Uint64Ptr returns a pointer to the provided value
func Uint64Ptr(v uint64) *uint64 { return &v }

// Float32Ptr returns a pointer to the provided value
func Float32Ptr(v float32) *float32 { return &v }

// Float64Ptr returns a pointer to the provided value
func Float64Ptr(v float64) *float64 { return &v }

// Complex64Ptr returns a pointer to the provided value
func Complex64Ptr(v complex64) *complex64 { return &v }

// Complex128Ptr returns a pointer to the provided value
func Complex128Ptr(v complex128) *complex128 { return &v }

// StringPtr returns a pointer to the provided value
func StringPtr(v string) *string { return &v }
