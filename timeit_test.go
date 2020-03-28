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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
