package metrics

import (
	"testing"

	"order-system/pkg/infra/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("metrics enabled", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Metrics.Enabled = true

		collector, err := New(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, collector)
	})

	t.Run("metrics disabled", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Metrics.Enabled = false

		collector, err := New(cfg)
		assert.Error(t, err)
		assert.Nil(t, collector)
	})
}

func TestRegister(t *testing.T) {
	cfg := &config.Config{}
	cfg.Metrics.Enabled = true
	collector, _ := New(cfg)

	t.Run("register new counter", func(t *testing.T) {
		err := collector.Register("test_counter", Counter, "test counter")
		assert.NoError(t, err)
	})

	t.Run("register duplicate metric", func(t *testing.T) {
		_ = collector.Register("test_counter", Counter, "initial registration")
		err := collector.Register("test_counter", Counter, "duplicate registration")
		assert.Error(t, err)
	})

	t.Run("register gauge", func(t *testing.T) {
		err := collector.Register("test_gauge", Gauge, "test gauge")
		assert.NoError(t, err)
	})

	t.Run("register histogram", func(t *testing.T) {
		err := collector.Register("test_histogram", Histogram, "test histogram")
		assert.NoError(t, err)
	})
}

func TestIncrementCounterAndGetCounter(t *testing.T) {
	cfg := &config.Config{}
	cfg.Metrics.Enabled = true
	collector, _ := New(cfg)

	labels := Labels{"label1": "value1"}

	t.Run("increment unregistered counter", func(t *testing.T) {
		collector.Register("unregistered", Gauge, "test gauge for error")
		collector.IncrementCounter("unregistered", 1.0, labels)
		value := collector.GetCounter("unregistered", labels)
		assert.Equal(t, float64(0), value)
	})

	t.Run("increment registered counter", func(t *testing.T) {
		collector.Register("test_counter", Counter, "test counter")

		collector.IncrementCounter("test_counter", 1.0, labels)
		value := collector.GetCounter("test_counter", labels)
		assert.Equal(t, float64(1), value)

		collector.IncrementCounter("test_counter", 2.0, labels)
		value = collector.GetCounter("test_counter", labels)
		assert.Equal(t, float64(3), value)
	})
}

func TestSetGaugeAndGetGauge(t *testing.T) {
	cfg := &config.Config{}
	cfg.Metrics.Enabled = true
	collector, _ := New(cfg)

	labels := Labels{"label1": "value1"}

	t.Run("set unregistered gauge", func(t *testing.T) {
		collector.Register("unregistered", Counter, "test counter for error")
		collector.SetGauge("unregistered", 1.0, labels)
		value := collector.GetGauge("unregistered", labels)
		assert.Equal(t, float64(0), value)
	})

	t.Run("set registered gauge", func(t *testing.T) {
		collector.Register("test_gauge", Gauge, "test gauge")

		collector.SetGauge("test_gauge", 1.0, labels)
		value := collector.GetGauge("test_gauge", labels)
		assert.Equal(t, float64(1), value)

		collector.SetGauge("test_gauge", 2.0, labels)
		value = collector.GetGauge("test_gauge", labels)
		assert.Equal(t, float64(2), value)
	})
}

func TestObserveHistogramAndGetHistogram(t *testing.T) {
	cfg := &config.Config{}
	cfg.Metrics.Enabled = true
	collector, _ := New(cfg)

	labels := Labels{"label1": "value1"}

	t.Run("observe unregistered histogram", func(t *testing.T) {
		collector.Register("unregistered", Gauge, "test gauge for error")
		collector.ObserveHistogram("unregistered", 1.0, labels)
		values := collector.GetHistogram("unregistered", labels)
		assert.Nil(t, values)
	})

	t.Run("observe registered histogram", func(t *testing.T) {
		collector.Register("test_histogram", Histogram, "test histogram")

		collector.ObserveHistogram("test_histogram", 1.0, labels)
		collector.ObserveHistogram("test_histogram", 2.0, labels)

		values := collector.GetHistogram("test_histogram", labels)
		assert.Equal(t, []float64{1.0, 2.0}, values)
	})
}
