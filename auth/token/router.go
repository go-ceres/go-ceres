//    Copyright 2022. Go-Ceres
//    Author https://github.com/go-ceres/go-ceres
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package token

type Router struct {
	trees       *node
	maxParams   uint16
	maxSections uint16
}

type RouterInfo struct {
	Permissions permissions
	Params      *Params
	Tsr         bool
	FullPath    string
}

func NewRouter() *Router {
	return &Router{
		trees:       new(node),
		maxParams:   0,
		maxSections: 0,
	}
}

func (r *Router) AddRouter(path string, permissions permissions) {
	r.trees.addRoute(path, permissions)
	// Update maxParams
	if paramsCount := countParams(path); paramsCount > r.maxParams {
		r.maxParams = paramsCount
	}

	if sectionsCount := countSections(path); sectionsCount > r.maxSections {
		r.maxSections = sectionsCount
	}
}

func (r *Router) GetRouter(path string, unescape bool) *RouterInfo {
	param := make(Params, 0, r.maxParams)
	skippedNodes := make([]skippedNode, 0, r.maxSections)
	value := r.trees.getValue(path, &param, &skippedNodes, unescape)
	return &RouterInfo{
		Permissions: value.permissions,
		Params:      value.params,
		Tsr:         value.tsr,
		FullPath:    value.fullPath,
	}
}
