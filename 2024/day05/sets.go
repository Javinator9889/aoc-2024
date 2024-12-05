// Mimics sets in Golang with O(m + n) time complexity
package main

// Performs the intersection between two slices
func intersection(a, b []int) (c []int) {
	// Verify length first to avoid unnecessary operations
	if len(a) == 0 || len(b) == 0 {
		return
	}
	m := make(map[int]bool)
	for _, item := range a {
		m[item] = true
	}
	for _, item := range b {
		if m[item] {
			c = append(c, item)
		}
	}
	c = _removeDups(c)
	return
}

func _removeDups(a []int) (b []int) {
	m := make(map[int]bool)
	for _, item := range a {
		if !m[item] {
			b = append(b, item)
			m[item] = true
		}
	}
	return
}
