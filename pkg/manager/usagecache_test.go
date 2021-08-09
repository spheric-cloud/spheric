/*
 * Copyright (c) 2021 by the OnMetal authors.
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

package manager

import (
	"github.com/onmetal/onmetal-api/pkg/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UsageCache", func() {

	var cache *UsageCache
	objectk10 := utils.MustParseObjectId("k1.g1/default/object0")
	objectk11 := utils.MustParseObjectId("k1.g1/default/object1")
	objectk12 := utils.MustParseObjectId("k1.g1/default/object2")
	//objectk20 := utils.MustParseObjectId("k2.g1/default/object0")
	//objectk21 := utils.MustParseObjectId("k2.g1/default/object1")

	BeforeEach(func() {
		cache = NewUsageCache(nil, nil)
	})

	Context("adding", func() {
		It("simple usage", func() {
			info := NewObjectUsageInfo("uses", objectk11)
			cache.replaceObjectUsageInfo(objectk10, info)
			Expect(cache.GetUsedObjectsFor(objectk10)).To(Equal(utils.NewObjectIds(objectk11)))
			Expect(cache.GetUsersFor(objectk11)).To(Equal(utils.NewObjectIds(objectk10)))
		})
		It("two relations", func() {
			info := NewObjectUsageInfo("uses", objectk11, "owner", objectk12)
			cache.replaceObjectUsageInfo(objectk10, info)
			Expect(cache.GetUsedObjectsFor(objectk10)).To(Equal(utils.NewObjectIds(objectk11, objectk12)))
			Expect(cache.GetUsersFor(objectk11)).To(Equal(utils.NewObjectIds(objectk10)))
			Expect(cache.GetUsersFor(objectk12)).To(Equal(utils.NewObjectIds(objectk10)))
		})
		It("query two relations", func() {
			info := NewObjectUsageInfo("uses", objectk11, "owner", objectk12)
			cache.replaceObjectUsageInfo(objectk10, info)
			Expect(cache.GetUsedObjectsForRelation(objectk10, "uses")).To(Equal(utils.NewObjectIds(objectk11)))
			Expect(cache.GetUsedObjectsForRelation(objectk10, "owner")).To(Equal(utils.NewObjectIds(objectk12)))
			Expect(cache.GetUsersForRelation(objectk11, "uses")).To(Equal(utils.NewObjectIds(objectk10)))
			Expect(cache.GetUsersForRelation(objectk11, "owner")).To(BeNil())
			Expect(cache.GetUsersForRelation(objectk12, "uses")).To(BeNil())
			Expect(cache.GetUsersForRelation(objectk12, "owner")).To(Equal(utils.NewObjectIds(objectk10)))
		})
	})
})
