package utils

import "fmt"

func Remove[T any](slice []T, s int) []T {
	if s >= len(slice) {
		fmt.Println("Index out of bounds", s, "for slice of length", len(slice), ".")
		return slice
	}
	return append(slice[:s], slice[s+1:]...)
}
