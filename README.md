Time Aggregation Library
========================

[![Tag](https://img.shields.io/github/tag/klmitch/timeit.svg)](https://github.com/klmitch/timeit/tags)
[![License](https://img.shields.io/hexpm/l/plug.svg)](https://github.com/klmitch/timeit/blob/master/LICENSE)
[![Test Report](https://travis-ci.org/klmitch/timeit.svg?branch=master)](https://travis-ci.org/klmitch/timeit)
[![Coverage Report](https://coveralls.io/repos/github/klmitch/timeit/badge.svg?branch=master)](https://coveralls.io/github/klmitch/timeit?branch=master)
[![Godoc](https://godoc.org/github.com/klmitch/timeit?status.svg)](http://godoc.org/github.com/klmitch/timeit)
[![Issue Tracker](https://img.shields.io/github/issues/klmitch/timeit.svg)](https://github.com/klmitch/timeit/issues)
[![Pull Request Tracker](https://img.shields.io/github/issues-pr/klmitch/timeit.svg)](https://github.com/klmitch/timeit/pulls)
[![Report Card](https://goreportcard.com/badge/github.com/klmitch/timeit)](https://goreportcard.com/report/github.com/klmitch/timeit)

The timeit repository contains an extremely simple library for collecting and computing statistical data regarding the amount of time it took for some function to be executed. To use, instantiate a pointer to an empty `Data` structure, then call `Data.TimeIt`, passing it the function to be called (or a function wrapping it). The `Data` structure will be updated each time `Data.TimeIt` is called; the `Data.Samples` element contains the total number of calls, `Data.Mean` contains the average time, `Data.Max` contains the maximum time, and `Data.Min` contains the minimum time. The variance and standard deviation may be calculated by calling the appropriate functions (`Data.Variance` and `Data.SampleVariance` for variance, `Data.StdDev` and `Data.SampleStdDev` for standard deviation). To manually update `Data` with a data point collected via some other means, use `Data.Update`.
