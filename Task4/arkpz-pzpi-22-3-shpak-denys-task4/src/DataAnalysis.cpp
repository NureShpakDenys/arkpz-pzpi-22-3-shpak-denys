#include "DataAnalysis.h"
#include <math.h>
#include <vector>

// Calculate the coefficient of variation of a dataset
float calculateCoefficientOfVariation(const std::vector<float>& data) {
  if (data.empty()) return 0.0;

  float sum = 0.0, mean, variance = 0.0, stddev;

  for (float value : data) {
    sum += value;
  }
  mean = sum / data.size();

  for (float value : data) {
    variance += pow(value - mean, 2);
  }
  variance /= data.size();
  stddev = sqrt(variance);

  return (mean != 0) ? stddev / mean : 0.0;
}