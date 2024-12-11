package cache

type Int2Func func(int, int) int

func Cached(f Int2Func) Int2Func {
	cache := make(map[int]map[int]int)
	return func(a, b int) int {
		if _, ok := cache[a]; !ok {
			cache[a] = make(map[int]int)
		}
		if v, ok := cache[a][b]; ok {
			return v
		}
		v := f(a, b)
		cache[a][b] = v
		return v
	}
}
