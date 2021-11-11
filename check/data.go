package check

type Data interface {
	String() string
	JsonString() (string, error)
}

type EmptyData struct {
}

func (EmptyData) String() string {
	return Na("")
}

func (EmptyData) JsonString() (string, error) {
	return "{}", nil
}

func Na(s string) string {
	if s == "" {
		return "n/a"
	}
	return s
}

func NonEmpty(strings ...string) []string {
	var ss []string
	for _, s := range strings {
		if s != "" {
			ss = append(ss, s)
		}
	}
	return ss
}
