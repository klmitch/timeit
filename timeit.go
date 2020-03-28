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
	"math"
	"time"
)

// Data contains the accumulated timing data.
type Data struct {
	Samples int64         // The number of samples developed so far
	Mean    time.Duration // The current running mean
	Max     time.Duration // Maximum sample seen so far
	Min     time.Duration // Minimum sample seen so far
	Next    *Data         // Another Data instance to update
	m2      time.Duration // Sum of square differences
}

// Update adds another sample to the Data structure.
func (d *Data) Update(sample time.Duration) {
	// Keep track of minimum and maximum
	if d.Samples == 0 || sample < d.Min {
		d.Min = sample
	}
	if d.Samples == 0 || sample > d.Max {
		d.Max = sample
	}

	// Update the sample count, mean, and m2 values
	d.Samples++
	delta1 := sample - d.Mean
	d.Mean = d.Mean + delta1/time.Duration(d.Samples)
	delta2 := sample - d.Mean
	d.m2 = d.m2 + delta1*delta2

	// Pass the sample on to Next
	if d.Next != nil {
		d.Next.Update(sample)
	}
}

// Variance returns the variance of the data.  This is the square of
// the standard deviation.  If no samples have been collected so far,
// this value will be 0.
func (d *Data) Variance() time.Duration {
	// Avoid divide by zero
	if d.Samples <= 0 {
		return time.Duration(0)
	}

	return d.m2 / time.Duration(d.Samples)
}

// SampleVariance returns the sample variance of the data.  This is
// the square of the standard deviation, and should probably be used
// in preference to Variance.  If only one sample has been collected
// so far, this value will be 0.
func (d *Data) SampleVariance() time.Duration {
	// Avoid divide by zero
	if d.Samples <= 1 {
		return time.Duration(0)
	}

	return d.m2 / time.Duration(d.Samples-1)
}

// StdDev returns the standard deviation of the data.  If no samples
// have been collected so far, this value will be 0.
func (d *Data) StdDev() time.Duration {
	return time.Duration(math.Sqrt(float64(d.Variance())))
}

// SampleStdDev returns the sample standard deviation of the data.
// This should probably be used in preference to StdDev.  If only one
// sample has been collected so far, this value will be 0.
func (d *Data) SampleStdDev() time.Duration {
	return time.Duration(math.Sqrt(float64(d.SampleVariance())))
}

// TimeIt runs a function and updates the data with the time it took
// for the function to execute.  It returns the time it took for the
// function to execute.
func (d *Data) TimeIt(fn func()) (delta time.Duration) {
	// Get the current time and arrange to update the data
	curr := time.Now()
	defer func() {
		delta = time.Since(curr)
		d.Update(delta)
	}()

	// Invoke the function
	fn()

	return
}
