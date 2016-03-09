package rtorrentexporter

import (
	"log"

	"github.com/mdlayher/rtorrent"
	"github.com/prometheus/client_golang/prometheus"
)

// A DownloadsCollector is a Prometheus collector for metrics regarding rTorrent
// downloads.
type DownloadsCollector struct {
	Downloads           *prometheus.Desc
	DownloadsStarted    *prometheus.Desc
	DownloadsStopped    *prometheus.Desc
	DownloadsComplete   *prometheus.Desc
	DownloadsIncomplete *prometheus.Desc
	DownloadsHashing    *prometheus.Desc
	DownloadsSeeding    *prometheus.Desc
	DownloadsLeeching   *prometheus.Desc
	DownloadsActive     *prometheus.Desc

	c *rtorrent.Client
}

// Verify that DownloadsCollector implements the prometheus.Collector interface.
var _ prometheus.Collector = &DownloadsCollector{}

// NewDownloadsCollector creates a new DownloadsCollector which collects metrics
// regarding rTorrent downloads.
func NewDownloadsCollector(c *rtorrent.Client) *DownloadsCollector {
	const (
		subsystem = "downloads"
	)

	return &DownloadsCollector{
		Downloads: prometheus.NewDesc(
			// Subsystem is used as name so we get "rtorrent_downloads"
			prometheus.BuildFQName(namespace, "", subsystem),
			"Total number of downloads.",
			nil,
			nil,
		),

		DownloadsStarted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "started"),
			"Number of started downloads.",
			nil,
			nil,
		),

		DownloadsStopped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "stopped"),
			"Number of stopped downloads.",
			nil,
			nil,
		),

		DownloadsComplete: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "complete"),
			"Number of complete downloads.",
			nil,
			nil,
		),

		DownloadsIncomplete: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "incomplete"),
			"Number of incomplete downloads.",
			nil,
			nil,
		),

		DownloadsHashing: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "hashing"),
			"Number of hashing downloads.",
			nil,
			nil,
		),

		DownloadsSeeding: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "seeding"),
			"Number of seeding downloads.",
			nil,
			nil,
		),

		DownloadsLeeching: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "leeching"),
			"Number of leeching downloads.",
			nil,
			nil,
		),

		DownloadsActive: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "active"),
			"Number of active downloads.",
			nil,
			nil,
		),

		c: c,
	}
}

// collect begins a metrics collection task for all metrics related to UniFi
// stations.
func (c *DownloadsCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	all, err := c.c.Downloads.All()
	if err != nil {
		return c.Downloads, err
	}

	started, err := c.c.Downloads.Started()
	if err != nil {
		return c.DownloadsStarted, err
	}

	stopped, err := c.c.Downloads.Stopped()
	if err != nil {
		return c.DownloadsStopped, err
	}

	complete, err := c.c.Downloads.Complete()
	if err != nil {
		return c.DownloadsComplete, err
	}

	incomplete, err := c.c.Downloads.Incomplete()
	if err != nil {
		return c.DownloadsIncomplete, err
	}

	hashing, err := c.c.Downloads.Hashing()
	if err != nil {
		return c.DownloadsHashing, err
	}

	seeding, err := c.c.Downloads.Seeding()
	if err != nil {
		return c.DownloadsSeeding, err
	}

	leeching, err := c.c.Downloads.Leeching()
	if err != nil {
		return c.DownloadsLeeching, err
	}

	active, err := c.c.Downloads.Active()
	if err != nil {
		return c.DownloadsActive, err
	}

	ch <- prometheus.MustNewConstMetric(
		c.Downloads,
		prometheus.GaugeValue,
		float64(len(all)),
	)

	ch <- prometheus.MustNewConstMetric(
		c.DownloadsStarted,
		prometheus.GaugeValue,
		float64(len(started)),
	)

	ch <- prometheus.MustNewConstMetric(
		c.DownloadsStopped,
		prometheus.GaugeValue,
		float64(len(stopped)),
	)

	ch <- prometheus.MustNewConstMetric(
		c.DownloadsComplete,
		prometheus.GaugeValue,
		float64(len(complete)),
	)

	ch <- prometheus.MustNewConstMetric(
		c.DownloadsIncomplete,
		prometheus.GaugeValue,
		float64(len(incomplete)),
	)

	ch <- prometheus.MustNewConstMetric(
		c.DownloadsHashing,
		prometheus.GaugeValue,
		float64(len(hashing)),
	)

	ch <- prometheus.MustNewConstMetric(
		c.DownloadsSeeding,
		prometheus.GaugeValue,
		float64(len(seeding)),
	)

	ch <- prometheus.MustNewConstMetric(
		c.DownloadsLeeching,
		prometheus.GaugeValue,
		float64(len(leeching)),
	)

	ch <- prometheus.MustNewConstMetric(
		c.DownloadsActive,
		prometheus.GaugeValue,
		float64(len(active)),
	)

	return nil, nil
}

// Describe sends the descriptors of each metric over to the provided channel.
// The corresponding metric values are sent separately.
func (c *DownloadsCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.Downloads,
		c.DownloadsStarted,
		c.DownloadsStopped,
		c.DownloadsComplete,
		c.DownloadsIncomplete,
		c.DownloadsHashing,
		c.DownloadsSeeding,
		c.DownloadsLeeching,
		c.DownloadsActive,
	}

	for _, d := range ds {
		ch <- d
	}
}

// Collect sends the metric values for each metric pertaining to the rTorrent
// downloads to the provided prometheus Metric channel.
func (c *DownloadsCollector) Collect(ch chan<- prometheus.Metric) {
	if desc, err := c.collect(ch); err != nil {
		log.Printf("[ERROR] failed collecting download metric %v: %v", desc, err)
		ch <- prometheus.NewInvalidMetric(desc, err)
		return
	}
}
