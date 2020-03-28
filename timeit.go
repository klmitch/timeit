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
	"math"
	"time"
)

// MarshalFlags contains a set of flags that controls how the Data
// will be marshaled into JSON or YAML.
type MarshalFlags uint8

// Recognized flags; these indicate which of the variance and standard
// deviation fields should be included in the marshaled object.
const (
	Variance       MarshalFlags = 1 << iota // Include Variance
	SampleVariance                          // Include SampleVariance
	StdDev                                  // Include StdDev
	SampleStdDev                            // Include SampleStdDev
)

// Data contains the accumulated timing data.
type Data struct {
	Samples int64         // The number of samples developed so far
	Mean    time.Duration // The current running mean
	Max     time.Duration // Maximum sample seen so far
	Min     time.Duration // Minimum sample seen so far
	Flags   MarshalFlags  // Bitmask of computed fields to marshal
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

// dataMarshaled contains the Data, along with the requested computed
// fields, which will then be marshaled into either JSON or YAML.
type dataMarshaled struct {
	Samples        *int64         `json:"samples" yaml:"samples"`
	Mean           *time.Duration `json:"mean" yaml:"mean"`
	Max            *time.Duration `json:"max" yaml:"max"`
	Min            *time.Duration `json:"min" yaml:"min"`
	Variance       *time.Duration `json:"variance,omitempty" yaml:"variance,omitempty"`
	SampleVariance *time.Duration `json:"sample_variance,omitempty" yaml:"sample_variance,omitempty"`
	StdDev         *time.Duration `json:"std_dev,omitempty" yaml:"std_dev,omitempty"`
	SampleStdDev   *time.Duration `json:"sample_std_dev,omitempty" yaml:"sample_std_dev,omitempty"`
}

// toData converts a dataMarshaled instance back into a Data instance.
// It guesses the Flags value based on the available data.
func (dm *dataMarshaled) toData(d *Data) {
	// Convert the basic data
	if dm.Samples != nil {
		d.Samples = *dm.Samples
	}
	if dm.Mean != nil {
		d.Mean = *dm.Mean
	}
	if dm.Max != nil {
		d.Max = *dm.Max
	}
	if dm.Min != nil {
		d.Min = *dm.Min
	}

	// Now handle the calculated values; go from the hardest to
	// recover m2 to the easiest, to attempt to be as accurate as
	// possible
	if dm.SampleStdDev != nil {
		d.Flags |= SampleStdDev
		if d.Samples > 1 {
			d.m2 = *dm.SampleStdDev * *dm.SampleStdDev * time.Duration(d.Samples-1)
		}
	}
	if dm.StdDev != nil {
		d.Flags |= StdDev
		d.m2 = *dm.StdDev * *dm.StdDev * time.Duration(d.Samples)
	}
	if dm.SampleVariance != nil {
		d.Flags |= SampleVariance
		if d.Samples > 1 {
			d.m2 = *dm.SampleVariance * time.Duration(d.Samples-1)
		}
	}
	if dm.Variance != nil {
		d.Flags |= Variance
		d.m2 = *dm.Variance * time.Duration(d.Samples)
	}
}

// marshaler constructs a dataMarshaled structure from Data.
func (d *Data) marshaler() *dataMarshaled {
	obj := &dataMarshaled{
		Samples: &d.Samples,
		Mean:    &d.Mean,
		Max:     &d.Max,
		Min:     &d.Min,
	}

	// Add requested computed fields
	if d.Flags == 0 || (d.Flags&Variance) != 0 {
		tmp := d.Variance()
		obj.Variance = &tmp
	}
	if d.Flags == 0 || (d.Flags&SampleVariance) != 0 {
		tmp := d.SampleVariance()
		obj.SampleVariance = &tmp
	}
	if d.Flags == 0 || (d.Flags&StdDev) != 0 {
		tmp := d.StdDev()
		obj.StdDev = &tmp
	}
	if d.Flags == 0 || (d.Flags&SampleStdDev) != 0 {
		tmp := d.SampleStdDev()
		obj.SampleStdDev = &tmp
	}

	return obj
}

// MarshalYAML implements yaml.Marshaler and allows a Data to be
// serialized intelligibly as YAML.
func (d *Data) MarshalYAML() (interface{}, error) {
	return d.marshaler(), nil
}

// UnmarshalYAML implements yaml.Unmarshaler and allows a Data to be
// deserialized intelligibly from YAML.  Note that round-tripping
// results in some inaccuracies in the calculations.
func (d *Data) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Unmarshal into a dataMarshaled struct
	dm := &dataMarshaled{}
	if err := unmarshal(dm); err != nil {
		return err
	}

	// Convert the dm to Data
	dm.toData(d)

	return nil
}

// MarshalJSON implements json.Marshaler and allows a Data to be
// serialized intelligibly as JSON.
func (d *Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.marshaler())
}

// UnmarshalJSON implements json.Unmarshaler and allows a Data to be
// deserialized intelligibly from JSON.  Note that round-tripping
// results in some inaccuracies in the calculations.
func (d *Data) UnmarshalJSON(text []byte) error {
	// Implement the noop convention
	if string(text) == "null" {
		return nil
	}

	// Unmarshal into a dataMarshaled struct
	dm := &dataMarshaled{}
	if err := json.Unmarshal(text, dm); err != nil {
		return err
	}

	// Convert the dm to Data
	dm.toData(d)

	return nil
}
