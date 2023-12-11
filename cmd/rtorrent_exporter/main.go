// Command rtorrent_exporter provides a Prometheus exporter for rTorrent.
package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/user"
	"runtime"

	"github.com/go-kit/log/level"

	"github.com/alecthomas/kingpin/v2"
	"github.com/mdlayher/rtorrent"
	rtorrentexporter "github.com/mdlayher/rtorrent_exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

func main() {
	var (
		toolkitFlags = kingpinflag.AddFlags(kingpin.CommandLine, ":9135")
		metricsPath  = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()

		rtorrentAddr     = kingpin.Flag("rtorrent.addr", "address of rTorrent XML-RPC server").Default("").String()
		rtorrentUsername = kingpin.Flag("rtorrent.username", "[optional] username used for HTTP Basic authentication with rTorrent XML-RPC server").Default("").String()
		rtorrentPassword = kingpin.Flag("rtorrent.password", "[optional] password used for HTTP Basic authentication with rTorrent XML-RPC server").Default("").String()
		maxProcs         = kingpin.Flag(
			"runtime.gomaxprocs", "The target number of CPUs Go will run on (GOMAXPROCS)",
		).Envar("GOMAXPROCS").Default("1").Int()
	)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("rtorrent_exporter"))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	if *rtorrentAddr == "" {
		level.Error(logger).Log("address of rTorrent XML-RPC server must be specified with '-rtorrent.addr' flag")
	}

	// Optionally enable HTTP Basic authentication
	var rt http.RoundTripper
	if u, p := *rtorrentUsername, *rtorrentPassword; u != "" && p != "" {
		rt = &authRoundTripper{
			Username: u,
			Password: p,
		}
	}

	c, err := rtorrent.New(*rtorrentAddr, rt)
	if err != nil {
		level.Error(logger).Log("cannot create rTorrent client: %v", err)
	}

	prometheus.MustRegister(rtorrentexporter.New(c))

	level.Info(logger).Log("starting rTorrent exporter on %q for server %q (authentication: %v)",
		*toolkitFlags, *rtorrentAddr, rt != nil)

	level.Info(logger).Log("msg", "Starting rtorrent_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())
	if user, err := user.Current(); err == nil && user.Uid == "0" {
		level.Warn(logger).Log("msg", "rTorrent Exporter is running as root user. This exporter is designed to run as unprivileged user, root is not required.")
	}
	runtime.GOMAXPROCS(*maxProcs)
	level.Debug(logger).Log("msg", "Go MAXPROCS", "procs", runtime.GOMAXPROCS(0))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	server := &http.Server{}
	if err := web.ListenAndServe(server, toolkitFlags, logger); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}

var _ http.RoundTripper = &authRoundTripper{}

// An authRoundTripper is a http.RoundTripper which adds HTTP Basic authentication
// to each HTTP request.
type authRoundTripper struct {
	Username string
	Password string
}

func (rt *authRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(rt.Username, rt.Password)
	return http.DefaultTransport.RoundTrip(r)
}
