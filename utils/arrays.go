package utils

import (
	"fmt"
)

func Remove[T any](slice []T, s int) []T {
	if s >= len(slice) {
		fmt.Println("Index out of bounds", s, "for slice of length", len(slice), ".")
		return slice
	}
	return append(slice[:s], slice[s+1:]...)
}

func RemoveElement[T comparable](slice []T, elem T) []T {
	for i, e := range slice {
		if e == elem {
			return Remove(slice, i)
		}
	}
	return slice
}

func Filter[T any](slice []T, f func(T) bool) []T {
	var filtered []T = make([]T, len(slice))
	n := 0
	for _, elem := range slice {
		if f(elem) {
			filtered[n] = elem
			n++
		}
	}
	return filtered[:n]
}

func Map[T any, U any](slice []T, f func(T) U) []U {
	var mapped []U = make([]U, len(slice))
	for i, elem := range slice {
		mapped[i] = f(elem)
	}
	return mapped
}

func Reduce[T any, U any](slice []T, f func(U, T) U, init U) U {
	var reduced U = init
	for _, elem := range slice {
		reduced = f(reduced, elem)
	}
	return reduced
}

func Any[T any](slice []T, f func(T) bool) bool {
	for _, elem := range slice {
		if f(elem) {
			return true
		}
	}
	return false
}

func All[T any](slice []T, f func(T) bool) bool {
	for _, elem := range slice {
		if !f(elem) {
			return false
		}
	}
	return true
}
