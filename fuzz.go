// +build ignore

package iso8601

func Fuzz(data []byte) int {
	_, err := Parse(data)
	if err == nil {
		return 1
	}
	return 0
}
