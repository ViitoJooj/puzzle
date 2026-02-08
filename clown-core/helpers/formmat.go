package helpers

import "strings"

func ReplaceProjectName(content []byte, name string) []byte {
	str := string(content)
	str = strings.ReplaceAll(str, "<projectname>", name)
	return []byte(str)
}
