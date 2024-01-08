package monitoring

import (
	"github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/ports"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var DefaultServerMonitor ports.WebsocketServerMonitor = newPrometheusWebsocketServerMonitor()

func init() {
	DefaultServerMonitor.Init()
}

type prometheusWebsocketServerMonitor struct {
	ConnectedShardsTotal prometheus.Gauge
}

func (p *prometheusWebsocketServerMonitor) Init() {
	prometheus.MustRegister(p.ConnectedShardsTotal)
}

func (p *prometheusWebsocketServerMonitor) RegisterClient(client ports.WebsocketServerClient) {
	p.ConnectedShardsTotal.Inc()
	log.Debugf("%T Registered client: %T", p, client)
}

func (p *prometheusWebsocketServerMonitor) UnregisterClient(client ports.WebsocketServerClient) {
	p.ConnectedShardsTotal.Dec()
	log.Debugf("%T Unregistered client: %T", p, client)
}

func newPrometheusWebsocketServerMonitor() *prometheusWebsocketServerMonitor {
	return &prometheusWebsocketServerMonitor{
		ConnectedShardsTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "connected_shards_total",
			Help: "Total number of connected shards",
		}),
	}
}
