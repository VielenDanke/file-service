package metrics

import (
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type cache struct {
	registerer prometheus.Registerer
	lock       sync.Mutex
	cVecs      map[string]*prometheus.CounterVec
	cs         map[string]prometheus.Counter
	gVecs      map[string]*prometheus.GaugeVec
	gs         map[string]prometheus.Gauge
	hVecs      map[string]*prometheus.HistogramVec
	hs         map[string]prometheus.Histogram
	sVecs      map[string]*prometheus.SummaryVec
	ss         map[string]prometheus.Summary
}

func newCache() *cache {
	return &cache{
		registerer: registerer,
		cVecs:      make(map[string]*prometheus.CounterVec),
		cs:         make(map[string]prometheus.Counter),
		gVecs:      make(map[string]*prometheus.GaugeVec),
		gs:         make(map[string]prometheus.Gauge),
		hVecs:      make(map[string]*prometheus.HistogramVec),
		hs:         make(map[string]prometheus.Histogram),
		sVecs:      make(map[string]*prometheus.SummaryVec),
		ss:         make(map[string]prometheus.Summary),
	}
}

func (c *cache) getCacheKey(namespace, subsystem, name string, labels []string) string {
	return strings.Join(append([]string{namespace, subsystem, name}, labels...), "||")
}

func (c *cache) getOrMakeCounterVec(opts prometheus.CounterOpts, labelNames []string) *prometheus.CounterVec {
	cacheKey := c.getCacheKey(opts.Namespace, opts.Subsystem, opts.Name, labelNames)
	c.lock.Lock()
	cv, cvExists := c.cVecs[cacheKey]
	if !cvExists {
		cv = prometheus.NewCounterVec(opts, labelNames)
		c.registerer.MustRegister(cv)
		c.cVecs[cacheKey] = cv
	}
	c.lock.Unlock()

	return cv
}

func (c *cache) getOrMakeCounter(opts prometheus.CounterOpts) prometheus.Counter {
	cacheKey := c.getCacheKey(opts.Namespace, opts.Subsystem, opts.Name, nil)
	c.lock.Lock()
	cn, cvExists := c.cs[cacheKey]
	if !cvExists {
		cn = prometheus.NewCounter(opts)
		c.registerer.MustRegister(cn)
		c.cs[cacheKey] = cn
	}
	c.lock.Unlock()

	return cn
}

func (c *cache) getOrMakeGaugeVec(opts prometheus.GaugeOpts, labelNames []string) *prometheus.GaugeVec {
	cacheKey := c.getCacheKey(opts.Namespace, opts.Subsystem, opts.Name, labelNames)
	c.lock.Lock()
	gv, gvExists := c.gVecs[cacheKey]
	if !gvExists {
		gv = prometheus.NewGaugeVec(opts, labelNames)
		c.registerer.MustRegister(gv)
		c.gVecs[cacheKey] = gv
	}
	c.lock.Unlock()

	return gv
}

func (c *cache) getOrMakeGauge(opts prometheus.GaugeOpts) prometheus.Gauge {
	cacheKey := c.getCacheKey(opts.Namespace, opts.Subsystem, opts.Name, nil)
	c.lock.Lock()
	g, gvExists := c.gs[cacheKey]
	if !gvExists {
		g = prometheus.NewGauge(opts)
		c.registerer.MustRegister(g)
		c.gs[cacheKey] = g
	}
	c.lock.Unlock()

	return g
}

func (c *cache) getOrMakeHistogramVec(opts prometheus.HistogramOpts, labelNames []string) *prometheus.HistogramVec {
	cacheKey := c.getCacheKey(opts.Namespace, opts.Subsystem, opts.Name, labelNames)
	c.lock.Lock()
	hv, hvExists := c.hVecs[cacheKey]
	if !hvExists {
		hv = prometheus.NewHistogramVec(opts, labelNames)
		c.registerer.MustRegister(hv)
		c.hVecs[cacheKey] = hv
	}
	c.lock.Unlock()

	return hv
}

func (c *cache) getOrMakeHistogram(opts prometheus.HistogramOpts) prometheus.Histogram {
	cacheKey := c.getCacheKey(opts.Namespace, opts.Subsystem, opts.Name, nil)
	c.lock.Lock()
	h, hvExists := c.hs[cacheKey]
	if !hvExists {
		h = prometheus.NewHistogram(opts)
		c.registerer.MustRegister(h)
		c.hs[cacheKey] = h
	}
	c.lock.Unlock()

	return h
}

func (c *cache) getOrMakeSummaryVec(opts prometheus.SummaryOpts, labelNames []string) *prometheus.SummaryVec {
	cacheKey := c.getCacheKey(opts.Namespace, opts.Subsystem, opts.Name, labelNames)
	c.lock.Lock()
	sv, hvExists := c.sVecs[cacheKey]
	if !hvExists {
		sv = prometheus.NewSummaryVec(opts, labelNames)
		c.registerer.MustRegister(sv)
		c.sVecs[cacheKey] = sv
	}
	c.lock.Unlock()

	return sv
}

func (c *cache) getOrMakeSummary(opts prometheus.SummaryOpts) prometheus.Summary {
	cacheKey := c.getCacheKey(opts.Namespace, opts.Subsystem, opts.Name, nil)
	c.lock.Lock()
	s, hvExists := c.ss[cacheKey]
	if !hvExists {
		s = prometheus.NewSummary(opts)
		c.registerer.MustRegister(s)
		c.ss[cacheKey] = s
	}
	c.lock.Unlock()

	return s
}

func GetOrMakeCounterVec(opts prometheus.CounterOpts, labelNames []string) *prometheus.CounterVec {
	return vc.getOrMakeCounterVec(opts, labelNames)
}

func GetOrMakeCounter(opts prometheus.CounterOpts) prometheus.Counter {
	return vc.getOrMakeCounter(opts)
}

func GetOrMakeGaugeVec(opts prometheus.GaugeOpts, labelNames []string) *prometheus.GaugeVec {
	return vc.getOrMakeGaugeVec(opts, labelNames)
}

func GetOrMakeGauge(opts prometheus.GaugeOpts) prometheus.Gauge {
	return vc.getOrMakeGauge(opts)
}

func GetOrMakeHistogramVec(opts prometheus.HistogramOpts, labelNames []string) *prometheus.HistogramVec {
	return vc.getOrMakeHistogramVec(opts, labelNames)
}

func GetOrMakeHistogram(opts prometheus.HistogramOpts) prometheus.Histogram {
	return vc.getOrMakeHistogram(opts)
}

func GetOrMakeSummaryVec(opts prometheus.SummaryOpts, labelNames []string) *prometheus.SummaryVec {
	return vc.getOrMakeSummaryVec(opts, labelNames)
}

func GetOrMakeSummary(opts prometheus.SummaryOpts) prometheus.Summary {
	return vc.getOrMakeSummary(opts)
}
