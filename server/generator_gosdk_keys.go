package main

import "strings"

func generateGoCode_keys(kindName string) string {
	return strings.Replace(`

`, "<ENTITY>", kindName, -1)
}
