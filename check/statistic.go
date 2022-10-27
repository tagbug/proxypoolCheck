package check

import "github.com/Sansui233/proxypool/pkg/proxy"

// Statistic for a proxy
type Stat struct {
	Speed    float64
	Delay    uint16
	ReqCount uint16
	Relay    bool
	Pool     bool
	OutIp    string
	Id       string
}

// Statistic array for proxies
type StatList []Stat

// ProxyStats stores proxies' statistics
var ProxyStats StatList

func init() {
	ProxyStats = make(StatList, 0)
}

// Update speed for a Stat
func (ps *Stat) UpdatePSSpeed(speed float64) {
	if ps.Speed < 60 && ps.Speed != 0 {
		ps.Speed = 0.3*ps.Speed + 0.7*speed
	} else {
		ps.Speed = speed
	}
}

// Update delay for a Stat
func (ps *Stat) UpdatePSDelay(delay uint16) {
	ps.Delay = delay
}

// Update out ip for a Stat
func (ps *Stat) UpdatePSOutIp(outIp string) {
	ps.OutIp = outIp
}

// Count + 1 for a Stat
func (ps *Stat) UpdatePSCount() {
	ps.ReqCount++
}

// Find a proxy's Stat in StatList
func (psList StatList) Find(p proxy.Proxy) (*Stat, bool) {
	s := p.Identifier()
	for i := range psList {
		if psList[i].Id == s {
			return &psList[i], true
		}
	}
	return nil, false
}

// Return proxies that request count more than a given nubmer
func (psList StatList) ReqCountThan(n uint16, pl []proxy.Proxy, reset bool) []proxy.Proxy {
	proxies := make([]proxy.Proxy, 0)
	for _, p := range pl {
		for j := range psList {
			if psList[j].ReqCount > n && p.Identifier() == psList[j].Id {
				proxies = append(proxies, p)
			}
		}
	}
	// reset request count
	if reset {
		for i := range psList {
			psList[i].ReqCount = 0
		}
	}
	return proxies
}
