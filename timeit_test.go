// Copyright (c) 2020 Kevin L. Mitchell
//
// Licensed under the Apache License, Version 2.0 (the "License"); you
// may not use this file except in compliance with the License.  You
// may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.  See the License for the specific language governing
// permissions and limitations under the License.

package timeit

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestDataUpdateBase(t *testing.T) {
	d := &Data{}

	d.Update(time.Duration(50))

	assert.Equal(t, &Data{
		Samples: 1,
		Mean:    time.Duration(50),
		Max:     time.Duration(50),
		Min:     time.Duration(50),
		m2:      time.Duration(0),
	}, d)
}

func TestDataUpdateNewMin(t *testing.T) {
	d := &Data{
		Samples: 1,
		Mean:    time.Duration(50),
		Max:     time.Duration(50),
		Min:     time.Duration(50),
		m2:      time.Duration(0),
	}

	d.Update(time.Duration(25))

	assert.Equal(t, &Data{
		Samples: 2,
		Mean:    time.Duration(38),
		Max:     time.Duration(50),
		Min:     time.Duration(25),
		m2:      time.Duration(325),
	}, d)
}

func TestDataUpdateNewMax(t *testing.T) {
	d := &Data{
		Samples: 1,
		Mean:    time.Duration(50),
		Max:     time.Duration(50),
		Min:     time.Duration(50),
		m2:      time.Duration(0),
	}

	d.Update(time.Duration(75))

	assert.Equal(t, &Data{
		Samples: 2,
		Mean:    time.Duration(62),
		Max:     time.Duration(75),
		Min:     time.Duration(50),
		m2:      time.Duration(325),
	}, d)
}

func TestDataUpdateNext(t *testing.T) {
	d := &Data{
		Next: &Data{},
	}

	d.Update(time.Duration(50))

	assert.Equal(t, &Data{
		Samples: 1,
		Mean:    time.Duration(50),
		Max:     time.Duration(50),
		Min:     time.Duration(50),
		Next: &Data{
			Samples: 1,
			Mean:    time.Duration(50),
			Max:     time.Duration(50),
			Min:     time.Duration(50),
			m2:      time.Duration(0),
		},
		m2: time.Duration(0),
	}, d)
}

func TestDataVarianceSamples0(t *testing.T) {
	d := &Data{
		m2: time.Duration(50),
	}

	result := d.Variance()

	assert.Equal(t, time.Duration(0), result)
}

func TestDataVarianceSamples1(t *testing.T) {
	d := &Data{
		Samples: 1,
		m2:      time.Duration(50),
	}

	result := d.Variance()

	assert.Equal(t, time.Duration(50), result)
}

func TestDataVarianceSamples2(t *testing.T) {
	d := &Data{
		Samples: 2,
		m2:      time.Duration(50),
	}

	result := d.Variance()

	assert.Equal(t, time.Duration(25), result)
}

func TestDataSampleVarianceSamples0(t *testing.T) {
	d := &Data{
		m2: time.Duration(50),
	}

	result := d.SampleVariance()

	assert.Equal(t, time.Duration(0), result)
}

func TestDataSampleVarianceSamples1(t *testing.T) {
	d := &Data{
		Samples: 1,
		m2:      time.Duration(50),
	}

	result := d.SampleVariance()

	assert.Equal(t, time.Duration(0), result)
}

func TestDataSampleVarianceSamples2(t *testing.T) {
	d := &Data{
		Samples: 2,
		m2:      time.Duration(50),
	}

	result := d.SampleVariance()

	assert.Equal(t, time.Duration(50), result)
}

func TestDataSampleVarianceSamples3(t *testing.T) {
	d := &Data{
		Samples: 3,
		m2:      time.Duration(50),
	}

	result := d.SampleVariance()

	assert.Equal(t, time.Duration(25), result)
}

func TestDataStdDev(t *testing.T) {
	d := &Data{
		Samples: 1,
		m2:      time.Duration(64),
	}

	result := d.StdDev()

	assert.Equal(t, time.Duration(8), result)
}

func TestDataSampleStdDev(t *testing.T) {
	d := &Data{
		Samples: 2,
		m2:      time.Duration(64),
	}

	result := d.SampleStdDev()

	assert.Equal(t, time.Duration(8), result)
}

func TestDataTimeIt(t *testing.T) {
	d := &Data{}

	result := d.TimeIt(func() { time.Sleep(100 * time.Millisecond) })

	assert.Equal(t, &Data{
		Samples: 1,
		Mean:    result,
		Max:     result,
		Min:     result,
		m2:      time.Duration(0),
	}, d)
}

func TestDataMarshaledToData(t *testing.T) {
	samples := int64(3)
	mean := time.Duration(50)
	max := time.Duration(75)
	min := time.Duration(25)
	variance := time.Duration(416)
	sampleVariance := time.Duration(625)
	stdDev := time.Duration(20)
	sampleStdDev := time.Duration(25)
	dm := &dataMarshaled{
		Samples:        &samples,
		Mean:           &mean,
		Max:            &max,
		Min:            &min,
		Variance:       &variance,
		SampleVariance: &sampleVariance,
		StdDev:         &stdDev,
		SampleStdDev:   &sampleStdDev,
	}
	result := &Data{}

	dm.toData(result)

	assert.Equal(t, &Data{
		Samples: 3,
		Mean:    time.Duration(50),
		Max:     time.Duration(75),
		Min:     time.Duration(25),
		Flags:   Variance | SampleVariance | StdDev | SampleStdDev,
		m2:      time.Duration(1248),
	}, result)
}

func TestDataMarshalerBase(t *testing.T) {
	d := &Data{
		Samples: 3,
		Mean:    time.Duration(50),
		Max:     time.Duration(75),
		Min:     time.Duration(25),
		m2:      time.Duration(1250),
	}

	result := d.marshaler()

	samples := int64(3)
	mean := time.Duration(50)
	max := time.Duration(75)
	min := time.Duration(25)
	variance := time.Duration(416)
	sampleVariance := time.Duration(625)
	stdDev := time.Duration(20)
	sampleStdDev := time.Duration(25)
	assert.Equal(t, &dataMarshaled{
		Samples:        &samples,
		Mean:           &mean,
		Max:            &max,
		Min:            &min,
		Variance:       &variance,
		SampleVariance: &sampleVariance,
		StdDev:         &stdDev,
		SampleStdDev:   &sampleStdDev,
	}, result)
}

func TestDataMarshalerVariance(t *testing.T) {
	d := &Data{
		Samples: 3,
		Mean:    time.Duration(50),
		Max:     time.Duration(75),
		Min:     time.Duration(25),
		Flags:   Variance,
		m2:      time.Duration(1250),
	}

	result := d.marshaler()

	samples := int64(3)
	mean := time.Duration(50)
	max := time.Duration(75)
	min := time.Duration(25)
	variance := time.Duration(416)
	assert.Equal(t, &dataMarshaled{
		Samples:  &samples,
		Mean:     &mean,
		Max:      &max,
		Min:      &min,
		Variance: &variance,
	}, result)
}

func TestDataMarshalerSampleVariance(t *testing.T) {
	d := &Data{
		Samples: 3,
		Mean:    time.Duration(50),
		Max:     time.Duration(75),
		Min:     time.Duration(25),
		Flags:   SampleVariance,
		m2:      time.Duration(1250),
	}

	result := d.marshaler()

	samples := int64(3)
	mean := time.Duration(50)
	max := time.Duration(75)
	min := time.Duration(25)
	sampleVariance := time.Duration(625)
	assert.Equal(t, &dataMarshaled{
		Samples:        &samples,
		Mean:           &mean,
		Max:            &max,
		Min:            &min,
		SampleVariance: &sampleVariance,
	}, result)
}

func TestDataMarshalerStdDev(t *testing.T) {
	d := &Data{
		Samples: 3,
		Mean:    time.Duration(50),
		Max:     time.Duration(75),
		Min:     time.Duration(25),
		Flags:   StdDev,
		m2:      time.Duration(1250),
	}

	result := d.marshaler()

	samples := int64(3)
	mean := time.Duration(50)
	max := time.Duration(75)
	min := time.Duration(25)
	stdDev := time.Duration(20)
	assert.Equal(t, &dataMarshaled{
		Samples: &samples,
		Mean:    &mean,
		Max:     &max,
		Min:     &min,
		StdDev:  &stdDev,
	}, result)
}

func TestDataMarshalerSampleStdDev(t *testing.T) {
	d := &Data{
		Samples: 3,
		Mean:    time.Duration(50),
		Max:     time.Duration(75),
		Min:     time.Duration(25),
		Flags:   SampleStdDev,
		m2:      time.Duration(1250),
	}

	result := d.marshaler()

	samples := int64(3)
	mean := time.Duration(50)
	max := time.Duration(75)
	min := time.Duration(25)
	sampleStdDev := time.Duration(25)
	assert.Equal(t, &dataMarshaled{
		Samples:      &samples,
		Mean:         &mean,
		Max:          &max,
		Min:          &min,
		SampleStdDev: &sampleStdDev,
	}, result)
}

func TestDataMarshalYAML(t *testing.T) {
	d := &Data{
		Samples: 3,
		Mean:    time.Duration(50),
		Max:     time.Duration(75),
		Min:     time.Duration(25),
		m2:      time.Duration(1250),
	}

	result, err := yaml.Marshal(d)

	require.NoError(t, err)
	actual := &dataMarshaled{}
	err = yaml.Unmarshal(result, actual)
	require.NoError(t, err)
	samples := int64(3)
	mean := time.Duration(50)
	max := time.Duration(75)
	min := time.Duration(25)
	variance := time.Duration(416)
	sampleVariance := time.Duration(625)
	stdDev := time.Duration(20)
	sampleStdDev := time.Duration(25)
	assert.Equal(t, &dataMarshaled{
		Samples:        &samples,
		Mean:           &mean,
		Max:            &max,
		Min:            &min,
		Variance:       &variance,
		SampleVariance: &sampleVariance,
		StdDev:         &stdDev,
		SampleStdDev:   &sampleStdDev,
	}, actual)
}

func TestDataUnmarshalYAMLBase(t *testing.T) {
	text := []byte(`---
samples: 3
mean: 50ns
max: 75ns
min: 25ns
variance: 416ns
sample_variance: 625ns
std_dev: 20ns
sample_std_dev: 25ns
`)
	result := &Data{}

	err := yaml.Unmarshal(text, result)

	assert.NoError(t, err)
	assert.Equal(t, &Data{
		Samples: 3,
		Mean:    time.Duration(50),
		Max:     time.Duration(75),
		Min:     time.Duration(25),
		Flags:   Variance | SampleVariance | StdDev | SampleStdDev,
		m2:      time.Duration(1248),
	}, result)
}

func TestDataUnmarshalYAMLError(t *testing.T) {
	unmarshal := func(data interface{}) error {
		return assert.AnError
	}
	result := &Data{}

	err := result.UnmarshalYAML(unmarshal)

	assert.Same(t, assert.AnError, err)
	assert.Equal(t, &Data{}, result)
}

func TestDataMarshalJSON(t *testing.T) {
	d := &Data{
		Samples: 3,
		Mean:    time.Duration(50),
		Max:     time.Duration(75),
		Min:     time.Duration(25),
		m2:      time.Duration(1250),
	}

	result, err := json.Marshal(d)

	require.NoError(t, err)
	actual := &dataMarshaled{}
	err = json.Unmarshal(result, actual)
	require.NoError(t, err)
	samples := int64(3)
	mean := time.Duration(50)
	max := time.Duration(75)
	min := time.Duration(25)
	variance := time.Duration(416)
	sampleVariance := time.Duration(625)
	stdDev := time.Duration(20)
	sampleStdDev := time.Duration(25)
	assert.Equal(t, &dataMarshaled{
		Samples:        &samples,
		Mean:           &mean,
		Max:            &max,
		Min:            &min,
		Variance:       &variance,
		SampleVariance: &sampleVariance,
		StdDev:         &stdDev,
		SampleStdDev:   &sampleStdDev,
	}, actual)
}

func TestDataUnmarshalJSONBase(t *testing.T) {
	text := []byte(`{
    "samples": 3,
    "mean": 50,
    "max": 75,
    "min": 25,
    "variance": 416,
    "sample_variance": 625,
    "std_dev": 20,
    "sample_std_dev": 25
}`)
	result := &Data{}

	err := json.Unmarshal(text, result)

	assert.NoError(t, err)
	assert.Equal(t, &Data{
		Samples: 3,
		Mean:    time.Duration(50),
		Max:     time.Duration(75),
		Min:     time.Duration(25),
		Flags:   Variance | SampleVariance | StdDev | SampleStdDev,
		m2:      time.Duration(1248),
	}, result)
}

func TestDataUnmarshalJSONNull(t *testing.T) {
	text := []byte(`null`)
	result := &Data{}

	err := json.Unmarshal(text, result)

	assert.NoError(t, err)
	assert.Equal(t, &Data{}, result)
}

func TestDataUnmarshalJSONError(t *testing.T) {
	text := []byte(`{"samples": "3"}`)
	result := &Data{}

	err := json.Unmarshal(text, result)

	assert.NotNil(t, err)
	assert.Equal(t, &Data{}, result)
}
