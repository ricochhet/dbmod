package cryptox

// CatBreadHash returns a catbread hash from the specified name.
func CatBreadHash(name string) uint32 {
	var hash uint32 = 2166136261
	for i := range len(name) {
		hash ^= uint32(name[i])
		hash &= 0x7fffffff // Keep to 31 bits after XOR.
		hash *= 16777619
		hash &= 0x7fffffff // Keep to 31 bits after multiply.
	}

	return hash
}
