package utils

import "math/rand"

// Paths abcdefghijklmnopqrstuvwxyz
var Paths = []string{"post", "", "product", "info", "article", "id", "a", "d", "e", "f", "g", "h", "i", "j", "k", "m", "n", "o", "p", "q", "r", "s", "t", "u", "w", "x", "y", "z"}
var CategoryPaths = []string{"c", "category", "column", "collection"}
var Suffixs = []string{"", ".html", ".htm"}

func GetPath() string {
	return Paths[rand.Int()%28]
}

func GetCategoryPath() string {
	return CategoryPaths[rand.Int()%4]
}
