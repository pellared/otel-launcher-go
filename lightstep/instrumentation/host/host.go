// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package host // import "github.com/lightstep/otel-launcher-go/lightstep/instrumentation/host"

import (
	"context"
	"fmt"
	"sync"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
)

// Host reports the work-in-progress conventional host metrics specified by OpenTelemetry.
type host struct {
	meter metric.Meter
}

// config contains optional settings for reporting host metrics.
type config struct {
	// MeterProvider sets the metric.MeterProvider.  If nil, the global
	// Provider will be used.
	MeterProvider metric.MeterProvider
}

// Option supports configuring optional settings for host metrics.
type Option interface {
	apply(*config)
}

// WithMeterProvider sets the Metric implementation to use for
// reporting.  If this option is not used, the global metric.MeterProvider
// will be used.  `provider` must be non-nil.
func WithMeterProvider(provider metric.MeterProvider) Option {
	return metricProviderOption{provider}
}

type metricProviderOption struct{ metric.MeterProvider }

func (o metricProviderOption) apply(c *config) {
	if o.MeterProvider != nil {
		c.MeterProvider = o.MeterProvider
	}
}

// Attribute sets.
var (
	// Attribute sets for CPU time measurements.

	AttributeCPUTimeUser   = []attribute.KeyValue{attribute.String("state", "user")}
	AttributeCPUTimeSystem = []attribute.KeyValue{attribute.String("state", "system")}
	AttributeCPUTimeOther  = []attribute.KeyValue{attribute.String("state", "other")}
	AttributeCPUTimeIdle   = []attribute.KeyValue{attribute.String("state", "idle")}

	// Attribute sets used for Memory measurements.

	AttributeMemoryAvailable = []attribute.KeyValue{attribute.String("state", "available")}
	AttributeMemoryUsed      = []attribute.KeyValue{attribute.String("state", "used")}

	// Attribute sets used for Network measurements.

	AttributeNetworkTransmit = []attribute.KeyValue{attribute.String("direction", "transmit")}
	AttributeNetworkReceive  = []attribute.KeyValue{attribute.String("direction", "receive")}
)

// newConfig computes a config from a list of Options.
func newConfig(opts ...Option) config {
	c := config{
		MeterProvider: global.MeterProvider(),
	}
	for _, opt := range opts {
		opt.apply(&c)
	}
	return c
}

// Start initializes reporting of host metrics using the supplied config.
func Start(opts ...Option) error {
	c := newConfig(opts...)
	if c.MeterProvider == nil {
		c.MeterProvider = global.MeterProvider()
	}
	h := newHost(c)
	return h.register()
}

func newHost(c config) *host {
	return &host{
		meter: c.MeterProvider.Meter("otel_launcher_go/host"),
	}
}

func (h *host) register() error {
	var (
		err error

		hostCPUTime           instrument.Float64ObservableCounter
		hostMemoryUsage       instrument.Int64ObservableUpDownCounter
		hostMemoryUtilization instrument.Float64ObservableGauge
		networkIOUsage        instrument.Int64ObservableCounter

		// lock prevents a race between batch observer and instrument registration.
		lock sync.Mutex
	)

	lock.Lock()
	defer lock.Unlock()

	if hostCPUTime, err = h.meter.Float64ObservableCounter(
		"system.cpu.time",
		instrument.WithUnit("s"),
		instrument.WithDescription(
			"Accumulated CPU time spent by this process host attributed by state (User, System, Other, Idle)",
		),
	); err != nil {
		return err
	}

	if hostMemoryUsage, err = h.meter.Int64ObservableUpDownCounter(
		"system.memory.usage",
		instrument.WithUnit("By"),
		instrument.WithDescription(
			"Memory usage of this process host attributed by memory state (Used, Available)",
		),
	); err != nil {
		return err
	}

	if hostMemoryUtilization, err = h.meter.Float64ObservableGauge(
		"system.memory.utilization",
		instrument.WithDescription(
			"Memory utilization of this process host attributed by memory state (Used, Available)",
		),
	); err != nil {
		return err
	}

	if networkIOUsage, err = h.meter.Int64ObservableCounter(
		"system.network.io",
		instrument.WithUnit("By"),
		instrument.WithDescription(
			"Bytes transferred attributed by direction (Transmit, Receive)",
		),
	); err != nil {
		return err
	}

	_, err = h.meter.RegisterCallback(
		func(ctx context.Context, obs metric.Observer) error {
			lock.Lock()
			defer lock.Unlock()

			hostTimeSlice, err := cpu.TimesWithContext(ctx, false)
			if err != nil {
				return err
			}
			if len(hostTimeSlice) != 1 {
				return fmt.Errorf("host CPU usage: incorrect summary count")
			}

			vmStats, err := mem.VirtualMemoryWithContext(ctx)
			if err != nil {
				return err
			}

			ioStats, err := net.IOCountersWithContext(ctx, false)
			if err != nil {
				return err
			}
			if len(ioStats) != 1 {
				return err
			}

			// Host CPU time
			hostTime := hostTimeSlice[0]
			obs.ObserveFloat64(hostCPUTime, hostTime.User, AttributeCPUTimeUser...)
			obs.ObserveFloat64(hostCPUTime, hostTime.System, AttributeCPUTimeSystem...)

			// Note: "other" is the sum of all other known states.
			other := hostTime.Nice +
				hostTime.Iowait +
				hostTime.Irq +
				hostTime.Softirq +
				hostTime.Steal +
				hostTime.Guest +
				hostTime.GuestNice

			obs.ObserveFloat64(hostCPUTime, other, AttributeCPUTimeOther...)
			obs.ObserveFloat64(hostCPUTime, hostTime.Idle, AttributeCPUTimeIdle...)

			// Host memory usage
			obs.ObserveInt64(hostMemoryUsage, int64(vmStats.Used), AttributeMemoryUsed...)
			obs.ObserveInt64(hostMemoryUsage, int64(vmStats.Available), AttributeMemoryAvailable...)

			// Host memory utilization
			obs.ObserveFloat64(hostMemoryUtilization, float64(vmStats.Used)/float64(vmStats.Total), AttributeMemoryUsed...)
			obs.ObserveFloat64(hostMemoryUtilization, float64(vmStats.Available)/float64(vmStats.Total), AttributeMemoryAvailable...)

			// Host network usage
			obs.ObserveInt64(networkIOUsage, int64(ioStats[0].BytesSent), AttributeNetworkTransmit...)
			obs.ObserveInt64(networkIOUsage, int64(ioStats[0].BytesRecv), AttributeNetworkReceive...)
			return nil
		},
		hostCPUTime,
		hostMemoryUsage,
		hostMemoryUtilization,
		networkIOUsage,
	)

	if err != nil {
		return err
	}

	return nil
}
