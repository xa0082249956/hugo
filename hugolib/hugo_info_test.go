// Copyright 2018 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"fmt"
	"testing"

	"github.com/xa0082249956/hugo/helpers"
	"github.com/stretchr/testify/require"
)

func TestHugoInfo(t *testing.T) {
	assert := require.New(t)

	assert.Equal(helpers.CurrentHugoVersion.Version(), hugoInfo.Version)
	assert.IsType(helpers.HugoVersionString(""), hugoInfo.Version)
	assert.Equal(CommitHash, hugoInfo.CommitHash)
	assert.Equal(BuildDate, hugoInfo.BuildDate)
	assert.Contains(hugoInfo.Generator, fmt.Sprintf("Hugo %s", hugoInfo.Version))

}
