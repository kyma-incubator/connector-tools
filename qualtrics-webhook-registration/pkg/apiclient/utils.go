package apiclient

func removeTrailingSlash(s string) string {
	if s[len(s)-1:len(s)] == "/" {
		return s[:len(s)-1]
	} else {
		return s
	}

}