# Package godsp

Package godsp is a Go package developed to support some basic signal processing functions using the discrete wavelet transform (DWT).

## Packages

- **godsp**: General functions on vectors or sets of vectors.
- **godsp/dbscan**: Implementation of DBSCAN (https://en.wikipedia.org/wiki/DBSCAN) to cluster histogram bins.
- **godsp/peaks**: Peak detection for time series
- **godsp/ppeaks**: Peak detection on the basis of persistent homology:
[https://www.sthu.org/blog/13-perstopology-peakdetection/index.html](https://www.sthu.org/blog/13-perstopology-peakdetection/index.html).
- **godsp/dwt**: Lifting implementation of the discrete wavelet transform using the Daubechies 4 wavelet. See:

  Ripples in Mathematics. The Discrete Wavelet Transform.  
   A. Jensen and A. la Cour-Harbo  
   Springer 2001  
   Section 3.4

## Installation

    $ go get github.com/goccmack/godsp
