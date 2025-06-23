package M

type StringEncoder string

func (s StringEncoder) Encode() ([]byte, error) {
	return []byte(s), nil
}

func (s StringEncoder) Length() int {
	return len(s)
}
