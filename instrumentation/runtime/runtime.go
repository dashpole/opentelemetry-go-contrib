// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	"context"
	"math"
	"runtime/metrics"
	"slices"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/semconv/v1.43.0/goconv"
)

// ScopeName is the instrumentation scope name.
const ScopeName = "go.opentelemetry.io/contrib/instrumentation/runtime"

const (
	goHeapMemory                = "/memory/classes/heap/stacks:bytes"
	goMemoryHeapFree            = "/memory/classes/heap/free:bytes"
	goMemoryHeapObjects         = "/memory/classes/heap/objects:bytes"
	goMemoryHeapUnused          = "/memory/classes/heap/unused:bytes"
	goMemoryMetadataMCacheFree  = "/memory/classes/metadata/mcache/free:bytes"
	goMemoryMetadataMCacheInuse = "/memory/classes/metadata/mcache/inuse:bytes"
	goMemoryMetadataMSpanFree   = "/memory/classes/metadata/mspan/free:bytes"
	goMemoryMetadataMSpanInuse  = "/memory/classes/metadata/mspan/inuse:bytes"
	goMemoryMetadataOther       = "/memory/classes/metadata/other:bytes"
	goMemoryOsStacks            = "/memory/classes/os-stacks:bytes"
	goMemoryOther               = "/memory/classes/other:bytes"
	goMemoryProfilingBuckets    = "/memory/classes/profiling/buckets:bytes"
	goMemoryLimit               = "/gc/gomemlimit:bytes"
	goMemoryAllocated           = "/gc/heap/allocs:bytes"
	goMemoryAllocations         = "/gc/heap/allocs:objects"
	goMemoryGoal                = "/gc/heap/goal:bytes"
	goGoroutines                = "/sched/goroutines:goroutines"
	goMaxProcs                  = "/sched/gomaxprocs:threads"
	goConfigGC                  = "/gc/gogc:percent"
	goSchedLatencies            = "/sched/latencies:seconds"
	goGCCycles                  = "/gc/cycles/total:gc-cycles"
	goGCPauses                  = "/sched/pauses/total/gc:seconds"
	goCPUUser                   = "/cpu/classes/user:cpu-seconds"
	goCPUGCMarkAssist           = "/cpu/classes/gc/mark/assist:cpu-seconds"
	goCPUGCMarkDedicated        = "/cpu/classes/gc/mark/dedicated:cpu-seconds"
	goCPUGCMarkIdle             = "/cpu/classes/gc/mark/idle:cpu-seconds"
	goCPUGCPause                = "/cpu/classes/gc/pause:cpu-seconds"
	goCPUScavengeAssist         = "/cpu/classes/scavenge/assist:cpu-seconds"
	goCPUScavengeBackground     = "/cpu/classes/scavenge/background:cpu-seconds"
	goCPUIdle                   = "/cpu/classes/idle:cpu-seconds"
)

// Start initializes reporting of runtime metrics using the supplied config.
// For goroutine scheduling and GC pause metrics, additionally see [NewProducer].
//
// All metrics emitted by Start are Development metrics in the Go runtime
// semantic conventions and are registered with StabilityDevelopment. The
// following metrics are produced when the MeterProvider is configured with
// WithMinimumStability(StabilityDevelopment) or when enabled via a View:
//
//	go.memory.used          By            Memory used by the Go runtime.
//	go.memory.limit         By            Go runtime memory limit configured by the user, if a limit exists.
//	go.memory.allocated     By            Memory allocated to the heap by the application.
//	go.memory.allocations   {allocation}  Count of allocations to the heap by the application.
//	go.memory.gc.goal       By            Heap size target for the end of the GC cycle.
//	go.goroutine.count      {goroutine}   Count of live goroutines.
//	go.processor.limit      {thread}      The number of OS threads that can execute user-level Go code simultaneously.
//	go.config.gogc          %             Heap size target percentage configured by the user, otherwise 100.
//
// The following metrics are also opt-in in the semantic conventions and are
// registered with WithOptIn. They are only produced when enabled via a View:
//
//	go.memory.gc.cycles     {gc_cycle}    Number of completed GC cycles.
//	go.cpu.time             s             Estimated CPU time spent by the Go runtime.
//
// The go.memory.type, go.memory.detailed_type, go.cpu.state, and
// go.cpu.detailed_state attributes are Development attributes. They are only
// recorded when the Meter reports them as enabled for StabilityDevelopment;
// otherwise, the values they distinguish are summed.
func Start(opts ...Option) error {
	c := newConfig(opts...)
	meter := c.MeterProvider.Meter(
		ScopeName,
		metric.WithInstrumentationVersion(Version),
		metric.WithSchemaURL(semconv.SchemaURL),
	)
	memoryUsed, err := goconv.NewMemoryUsed(
		meter,
		metricx.WithDefaultAttributes((goconv.MemoryUsed{}).AttrMemoryType("").Key),
		metricx.WithStability(metricx.StabilityDevelopment),
	)
	if err != nil {
		return err
	}
	memoryLimit, err := goconv.NewMemoryLimit(
		meter,
		metricx.WithStability(metricx.StabilityDevelopment),
	)
	if err != nil {
		return err
	}
	memoryAllocated, err := goconv.NewMemoryAllocated(
		meter,
		metricx.WithStability(metricx.StabilityDevelopment),
	)
	if err != nil {
		return err
	}
	memoryAllocations, err := goconv.NewMemoryAllocations(
		meter,
		metricx.WithStability(metricx.StabilityDevelopment),
	)
	if err != nil {
		return err
	}
	memoryGCGoal, err := goconv.NewMemoryGCGoal(
		meter,
		metricx.WithStability(metricx.StabilityDevelopment),
	)
	if err != nil {
		return err
	}
	goroutineCount, err := goconv.NewGoroutineCount(
		meter,
		metricx.WithStability(metricx.StabilityDevelopment),
	)
	if err != nil {
		return err
	}
	processorLimit, err := goconv.NewProcessorLimit(
		meter,
		metricx.WithStability(metricx.StabilityDevelopment),
	)
	if err != nil {
		return err
	}
	configGogc, err := goconv.NewConfigGogc(
		meter,
		metricx.WithStability(metricx.StabilityDevelopment),
	)
	if err != nil {
		return err
	}
	gcCycles, err := goconv.NewMemoryGCCyclesObservable(
		meter,
		metricx.WithStability(metricx.StabilityDevelopment),
		metricx.WithOptIn(),
	)
	if err != nil {
		return err
	}
	cpuTime, err := goconv.NewCPUTimeObservable(
		meter,
		metricx.WithDefaultAttributes((goconv.CPUTimeObservable{}).AttrCPUState("").Key),
		metricx.WithStability(metricx.StabilityDevelopment),
		metricx.WithOptIn(),
	)
	if err != nil {
		return err
	}

	memoryClasses := newMemoryClasses(memoryUsed)
	cpuStates := newCPUStates(cpuTime)
	memoryObservations := newClassObservations(memoryClasses)
	cpuObservations := newClassObservations(cpuStates)
	instruments := []metric.Observable{
		memoryUsed.Inst(),
		memoryLimit.Inst(),
		memoryAllocated.Inst(),
		memoryAllocations.Inst(),
		memoryGCGoal.Inst(),
		goroutineCount.Inst(),
		processorLimit.Inst(),
		configGogc.Inst(),
		gcCycles.Inst(),
		cpuTime.Inst(),
	}
	metricNames := append(slices.Clone(runtimeMetrics), goGCCycles)
	for _, state := range cpuStates {
		metricNames = append(metricNames, state.runtimeMetric)
	}

	collector := newCollector(c.MinimumRefreshInterval, metricNames)
	var lock sync.Mutex
	_, err = meter.RegisterCallback(
		func(ctx context.Context, o metric.Observer) error {
			lock.Lock()
			defer lock.Unlock()
			collector.refresh()
			for _, obs := range memoryObservations.enabled(ctx, meter) {
				var v int64
				for _, name := range obs.runtimeMetrics {
					v += collector.getInt(name)
				}
				o.ObserveInt64(memoryUsed.Inst(), v, obs.attributes)
			}
			// Only observe the limit metric if a limit exists
			if limit := collector.getInt(goMemoryLimit); limit != math.MaxInt64 {
				o.ObserveInt64(memoryLimit.Inst(), limit)
			}
			o.ObserveInt64(memoryAllocated.Inst(), collector.getInt(goMemoryAllocated))
			o.ObserveInt64(memoryAllocations.Inst(), collector.getInt(goMemoryAllocations))
			o.ObserveInt64(memoryGCGoal.Inst(), collector.getInt(goMemoryGoal))
			o.ObserveInt64(goroutineCount.Inst(), collector.getInt(goGoroutines))
			o.ObserveInt64(processorLimit.Inst(), collector.getInt(goMaxProcs))
			o.ObserveInt64(configGogc.Inst(), collector.getInt(goConfigGC))
			o.ObserveInt64(gcCycles.Inst(), collector.getInt(goGCCycles))
			for _, obs := range cpuObservations.enabled(ctx, meter) {
				var v float64
				for _, name := range obs.runtimeMetrics {
					v += collector.getFloat(name)
				}
				o.ObserveFloat64(cpuTime.Inst(), v, obs.attributes)
			}
			return nil
		},
		instruments...,
	)
	if err != nil {
		return err
	}
	return nil
}

// These are the metrics we actually fetch from the go runtime.
var runtimeMetrics = []string{
	goHeapMemory,
	goMemoryHeapFree,
	goMemoryHeapObjects,
	goMemoryHeapUnused,
	goMemoryMetadataMCacheFree,
	goMemoryMetadataMCacheInuse,
	goMemoryMetadataMSpanFree,
	goMemoryMetadataMSpanInuse,
	goMemoryMetadataOther,
	goMemoryOsStacks,
	goMemoryOther,
	goMemoryProfilingBuckets,
	goMemoryLimit,
	goMemoryAllocated,
	goMemoryAllocations,
	goMemoryGoal,
	goGoroutines,
	goMaxProcs,
	goConfigGC,
}

// attrClass maps a runtime/metrics sample to the general (e.g. go.memory.type)
// and detailed (e.g. go.memory.detailed_type) Development attributes that
// describe it.
type attrClass struct {
	runtimeMetric     string
	general, detailed attribute.KeyValue
}

func newMemoryClasses(memoryUsed goconv.MemoryUsed) []attrClass {
	newClass := func(runtimeMetric string, mType goconv.MemoryTypeAttr, detailedType string) attrClass {
		return attrClass{
			runtimeMetric: runtimeMetric,
			general:       memoryUsed.AttrMemoryType(mType),
			detailed:      memoryUsed.AttrMemoryDetailedType(detailedType),
		}
	}
	return []attrClass{
		newClass(goHeapMemory, goconv.MemoryTypeStack, "heap/stacks"),
		newClass(goMemoryHeapFree, goconv.MemoryTypeOther, "heap/free"),
		newClass(goMemoryHeapObjects, goconv.MemoryTypeOther, "heap/objects"),
		newClass(goMemoryHeapUnused, goconv.MemoryTypeOther, "heap/unused"),
		newClass(goMemoryMetadataMCacheFree, goconv.MemoryTypeOther, "metadata/mcache/free"),
		newClass(goMemoryMetadataMCacheInuse, goconv.MemoryTypeOther, "metadata/mcache/inuse"),
		newClass(goMemoryMetadataMSpanFree, goconv.MemoryTypeOther, "metadata/mspan/free"),
		newClass(goMemoryMetadataMSpanInuse, goconv.MemoryTypeOther, "metadata/mspan/inuse"),
		newClass(goMemoryMetadataOther, goconv.MemoryTypeOther, "metadata/other"),
		newClass(goMemoryOsStacks, goconv.MemoryTypeOther, "os-stacks"),
		newClass(goMemoryOther, goconv.MemoryTypeOther, "other"),
		newClass(goMemoryProfilingBuckets, goconv.MemoryTypeOther, "profiling/buckets"),
	}
}

func newCPUStates(cpuTime goconv.CPUTimeObservable) []attrClass {
	newState := func(runtimeMetric string, state goconv.CPUStateAttr, detailedState string) attrClass {
		return attrClass{
			runtimeMetric: runtimeMetric,
			general:       cpuTime.AttrCPUState(state),
			detailed:      cpuTime.AttrCPUDetailedState(detailedState),
		}
	}
	return []attrClass{
		newState(goCPUUser, goconv.CPUStateUser, "user"),
		newState(goCPUGCMarkAssist, goconv.CPUStateGC, "gc/mark/assist"),
		newState(goCPUGCMarkDedicated, goconv.CPUStateGC, "gc/mark/dedicated"),
		newState(goCPUGCMarkIdle, goconv.CPUStateGC, "gc/mark/idle"),
		newState(goCPUGCPause, goconv.CPUStateGC, "gc/pause"),
		newState(goCPUScavengeAssist, goconv.CPUStateScavenge, "scavenge/assist"),
		newState(goCPUScavengeBackground, goconv.CPUStateScavenge, "scavenge/background"),
		newState(goCPUIdle, goconv.CPUStateIdle, "idle"),
	}
}

// observation is the sum of runtime/metrics samples observed with a single
// attribute set.
type observation struct {
	runtimeMetrics []string
	attributes     metric.ObserveOption
}

// classObservations holds the observations of a set of attrClasses for each
// combination of enabled general and detailed Development attributes. They
// are precomputed so that collection does not allocate attribute sets.
type classObservations struct {
	generalKey, detailedKey attribute.Key
	// byEnabled is indexed by [general enabled][detailed enabled].
	byEnabled [2][2][]observation
}

func newClassObservations(classes []attrClass) classObservations {
	c := classObservations{
		generalKey:  classes[0].general.Key,
		detailedKey: classes[0].detailed.Key,
	}
	for general := range 2 {
		for detailed := range 2 {
			obs := &c.byEnabled[general][detailed]
			index := make(map[attribute.Distinct]int)
			for _, class := range classes {
				var kvs []attribute.KeyValue
				if general == 1 {
					kvs = append(kvs, class.general)
				}
				if detailed == 1 {
					kvs = append(kvs, class.detailed)
				}
				set := attribute.NewSet(kvs...)
				if i, ok := index[set.Equivalent()]; ok {
					(*obs)[i].runtimeMetrics = append((*obs)[i].runtimeMetrics, class.runtimeMetric)
					continue
				}
				index[set.Equivalent()] = len(*obs)
				*obs = append(*obs, observation{
					runtimeMetrics: []string{class.runtimeMetric},
					attributes:     metric.WithAttributeSet(set),
				})
			}
		}
	}
	return c
}

// enabled returns the observations to record given the Development
// attributes m reports as enabled.
func (c classObservations) enabled(ctx context.Context, m metric.Meter) []observation {
	var general, detailed int
	if metricx.AttributeEnabled(ctx, m, c.generalKey, metricx.StabilityDevelopment) {
		general = 1
	}
	if metricx.AttributeEnabled(ctx, m, c.detailedKey, metricx.StabilityDevelopment) {
		detailed = 1
	}
	return c.byEnabled[general][detailed]
}

type goCollector struct {
	// now is used to replace the implementation of time.Now for testing
	now func() time.Time
	// lastCollect tracks the last time metrics were refreshed
	lastCollect time.Time
	// minimumInterval is the minimum amount of time between calls to metrics.Read
	minimumInterval time.Duration
	// sampleBuffer is populated by runtime/metrics
	sampleBuffer []metrics.Sample
	// sampleMap allows us to easily get the value of a single metric
	sampleMap map[string]*metrics.Sample
}

func newCollector(minimumInterval time.Duration, metricNames []string) *goCollector {
	g := &goCollector{
		sampleBuffer:    make([]metrics.Sample, 0, len(metricNames)),
		sampleMap:       make(map[string]*metrics.Sample, len(metricNames)),
		minimumInterval: minimumInterval,
		now:             time.Now,
	}
	for _, metricName := range metricNames {
		g.sampleBuffer = append(g.sampleBuffer, metrics.Sample{Name: metricName})
		// sampleMap references a position in the sampleBuffer slice. If an
		// element is appended to sampleBuffer, it must be added to sampleMap
		// for the sample to be accessible in sampleMap.
		g.sampleMap[metricName] = &g.sampleBuffer[len(g.sampleBuffer)-1]
	}
	return g
}

func (g *goCollector) refresh() {
	now := g.now()
	if now.Sub(g.lastCollect) < g.minimumInterval {
		// refresh was invoked more frequently than allowed by the minimum
		// interval. Do nothing.
		return
	}
	metrics.Read(g.sampleBuffer)
	g.lastCollect = now
}

func (g *goCollector) getInt(name string) int64 {
	if s, ok := g.sampleMap[name]; ok && s.Value.Kind() == metrics.KindUint64 {
		v := s.Value.Uint64()
		if v > math.MaxInt64 {
			return math.MaxInt64
		}
		return int64(v)
	}
	return 0
}

func (g *goCollector) getFloat(name string) float64 {
	if s, ok := g.sampleMap[name]; ok && s.Value.Kind() == metrics.KindFloat64 {
		return s.Value.Float64()
	}
	return 0
}

func (g *goCollector) getHistogram(name string) *metrics.Float64Histogram {
	if s, ok := g.sampleMap[name]; ok && s.Value.Kind() == metrics.KindFloat64Histogram {
		return s.Value.Float64Histogram()
	}
	return nil
}
