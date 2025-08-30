// Package utils contains utility functions
package utils

func RemoveIndex[T any](slice []T, index int) []T {
	slice = append(slice[:index], slice[index+1:]...)
	return slice
}
