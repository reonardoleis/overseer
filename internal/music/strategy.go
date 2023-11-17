package music

var (
	_ MusicStrategy = (*Youtube)(nil)
)

type Provider int

const (
	YOUTUBE Provider = iota

	DEFAULT = YOUTUBE
)

type MusicStrategy interface {
	Query(query string) (string, error)
	GetAudio(url string) ([]byte, error)
}

func New(provider ...Provider) MusicStrategy {
	p := YOUTUBE
	if len(provider) != 0 {
		p = provider[0]
	}

	switch p {
	case YOUTUBE:
		return &Youtube{}
	default:
		return New(DEFAULT)
	}
}
