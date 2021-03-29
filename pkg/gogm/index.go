// Copyright (c) 2020 MindStand Technologies, Inc
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package gogm

import (
	"github.com/cornelk/hashmap"
)

//drops all known indexes
func dropAllIndexesAndConstraints() error {
	if neoVersion >= 4 {
		return dropAllIndexesAndConstraintsV4()
	}

	return dropAllIndexesAndConstraintsV3()
}

//creates all indexes
func createAllIndexesAndConstraints(mappedTypes *hashmap.HashMap) error {
	if neoVersion >= 4 {
		return createAllIndexesAndConstraintsV4(mappedTypes)
	}

	return createAllIndexesAndConstraintsV3(mappedTypes)
}

//verifies all indexes
func verifyAllIndexesAndConstraints(mappedTypes *hashmap.HashMap) error {
	if neoVersion >= 4 {
		return verifyAllIndexesAndConstraintsV4(mappedTypes)
	}

	return verifyAllIndexesAndConstraintsV3(mappedTypes)
}
