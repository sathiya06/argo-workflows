package telemetry

import (
	"testing"

	"github.com/argoproj/argo-workflows/v3/util/logging"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

func TestViewDisable(t *testing.T) {
	// Same metric as TestMetrics, but disabled by a view
	ctx := logging.TestContext(t.Context())
	m, te, err := createTestMetrics(ctx, &Config{
		Modifiers: map[string]Modifier{
			nameTestingHistogram: {
				Disabled: true,
			},
		},
	})
	require.NoError(t, err)
	m.TestingHistogramRecord(ctx, 5)
	attribs := attribute.NewSet()
	_, err = te.GetFloat64HistogramData(ctx, nameTestingHistogram, &attribs)
	require.Error(t, err)
}

func TestViewDisabledAttributes(t *testing.T) {
	ctx := logging.TestContext(t.Context())
	// Disable the error cause attribute
	m, te, err := createTestMetrics(ctx, &Config{
		Modifiers: map[string]Modifier{
			nameTestingCounter: {
				DisabledAttributes: []string{AttribErrorCause},
			},
		},
	})
	require.NoError(t, err)
	// Submit a couple of errors
	m.TestingErrorA(ctx)
	m.TestingErrorB(ctx)
	// See if we can find this with the attributes, we should not be able to
	attribsFail := attribute.NewSet(attribute.String(AttribErrorCause, string(errorCauseTestingA)))
	_, err = te.GetInt64CounterValue(ctx, nameTestingCounter, &attribsFail)
	require.Error(t, err)
	// Find a sum of all error types
	attribsSuccess := attribute.NewSet()
	val, err := te.GetInt64CounterValue(ctx, nameTestingCounter, &attribsSuccess)
	require.NoError(t, err)
	// Sum of the two submitted errors is 2
	assert.Equal(t, int64(2), val)
}

func TestViewHistogramBuckets(t *testing.T) {
	// Same metric as TestMetrics, but buckets changed
	bounds := []float64{1.0, 3.0, 5.0, 10.0}
	ctx := logging.TestContext(t.Context())
	m, te, err := createTestMetrics(ctx, &Config{
		Modifiers: map[string]Modifier{
			nameTestingHistogram: {
				HistogramBuckets: bounds,
			},
		},
	})
	require.NoError(t, err)
	m.TestingHistogramRecord(ctx, 5)
	attribs := attribute.NewSet()
	val, err := te.GetFloat64HistogramData(ctx, nameTestingHistogram, &attribs)
	require.NoError(t, err)
	assert.Equal(t, bounds, val.Bounds)
	assert.Equal(t, []uint64{0, 0, 1, 0, 0}, val.BucketCounts)
}
