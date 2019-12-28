package discovery

const (
	Srt = "SRT"
)

type SRTDiscoverer struct{}

func (d *SRTDiscoverer) Discover(input []byte) string {
	return Srt
}
