// Copyright 2023 Specter Ops, Inc.
//
// Licensed under the Apache License, Version 2.0
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
//
// SPDX-License-Identifier: Apache-2.0

import * as types from './types';

export const setEntityInfoOpen = (open: boolean): types.EntityInfoActionTypes => {
    return {
        type: types.ENTITY_INFO_OPEN,
        open: open,
    };
};

export const setSelectedNode = (selectedNode: types.SelectedNode): types.EntityInfoActionTypes => {
    return {
        type: types.SET_SELECTED_NODE,
        selectedNode,
    };
};

export const addExpandedRelationship = (payload: string) => {
    return {
        type: types.ADD_EXPANDED_RELATIONSHIP,
        payload,
    };
};
export const removeExpandedRelationship = (payload: string) => {
    return {
        type: types.REMOVE_EXPANDED_RELATIONSHIP,
        payload,
    };
};
export const setExpandedRelationship = (payload: string[]) => {
    return {
        type: types.SET_EXPANDED_RELATIONSHIP,
        payload,
    };
};
export const clearExpandedRelationship = () => {
    return {
        type: types.CLEAR_EXPANDED_RELATIONSHIP,
    };
};
