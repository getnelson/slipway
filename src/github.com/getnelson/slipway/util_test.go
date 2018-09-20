//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//:
//:   Licensed under the Apache License, Version 2.0 (the "License");
//:   you may not use this file except in compliance with the License.
//:   You may obtain a copy of the License at
//:
//:       http://www.apache.org/licenses/LICENSE-2.0
//:
//:   Unless required by applicable law or agreed to in writing, software
//:   distributed under the License is distributed on an "AS IS" BASIS,
//:   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//:   See the License for the specific language governing permissions and
//:   limitations under the License.
//:
//: ----------------------------------------------------------------------------
package main

import (
	"testing"
)

func checkUnitNameFromImage(image string, oImage string, oTag string, t *testing.T) {
	e1 := oImage
	e2 := oTag
	last, tag := getUnitNameFromDockerContainer(image)

	if tag != e2 {
		t.Error(tag, e2)
	}
	if last != e1 {
		t.Error(last, e1)
	}
}

func TestUnitNameExtractor(t *testing.T) {
	checkUnitNameFromImage("docker.yourcompany.com/units/aloha-1.0:1.0.9", "aloha", "1.0.9", t)
	checkUnitNameFromImage("docker.yourcompany.com/units/aloha:1.0.9", "aloha", "1.0.9", t)
}
