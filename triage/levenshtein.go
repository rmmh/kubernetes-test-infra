/*
 * Copyright 2017 The Kubernetes Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package editdistance

// Compute the edit distance between two strings using the Levenshtein algorithm.
// If limit is given, exit early when the edit distance is guaranteed to be greater.
func LevenshteinDistance(a, b string, limit int) int {
	if a == b {
		return 0
	}

	if len(a) < len(b) {
		a, b = b, a
	}

	// create two work vectors of integer differences
	v0 := make([]int, len(b)+1)
	v1 := make([]int, len(v0))

	// initialize v0 (the previous row of distances)
	// this row is A[0][i]: edit distance for an empty a
	// the distance is just the number of characters to delete from b
	// note: v1 becomes v0 on loop entry
	for i := 0; i < len(v0); i++ {
		v1[i] = i
	}

	for i := 0; i < len(a); i++ {
		// calculate v1 (current row distances) from the previous row v0

		// swap v0, v1 for iteration
		v0, v1 = v1, v0

		// first element of v1 is A[i+1][0]
		//   edit distance is delete (i+1) chars from s to match empty t
		v1[0] = i + 1

		// use formula to fill in the rest of the row
		for j := 0; j < len(b); j++ {
			diag := 0
			if a[i] != b[j] {
				diag = 1
			}
			v1[j+1] = min(v1[j]+1, min(v0[j+1]+1, v0[j]+diag))
		}
	}

	return v1[len(b)]
}
