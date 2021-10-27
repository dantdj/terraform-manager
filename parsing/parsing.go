package parsing

import (
	"bufio"
	"os"
	"strings"
)

func ParseShaList(file *os.File) map[string]string {
	shaMap := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), "  ")
		shaMap[row[1]] = row[0]
	}
	return shaMap
}
