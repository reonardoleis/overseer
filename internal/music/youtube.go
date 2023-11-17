package music

type Youtube struct {
}

func (y Youtube) Query(query string) (string, error) {
	return "", nil
}

func (y Youtube) GetAudio(url string) ([]byte, error) {
	return nil, nil
}
