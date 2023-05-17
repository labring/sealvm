// Copyright © 2022 sealos.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package template

import (
	"bytes"
	"fmt"
	"testing"
)

func TestTemplateSemverCompare(t *testing.T) {
	v, b, e := TryParse(`
echo "{{ $result := readFile .PublicKey }}{{ b64enc $result }}"
echo "{{ $result := readFile .PrivateKey }}{{ b64enc $result }}"
`)
	if e != nil {
		t.Errorf("parse err: %v", e)
	}
	if !b {
		t.Errorf("parse failed: %v", b)
	}

	values := map[string]string{
		"PublicKey":  "/Users/cuisongliu/.ssh/id_rsa.pub",
		"PrivateKey": "/Users/cuisongliu/.ssh/id_rsa",
	}

	out := bytes.NewBuffer(nil)
	execErr := v.Execute(out, values)
	if execErr != nil {
		t.Errorf("template exec err: %v", execErr)
	}

	fmt.Println(out)
}
