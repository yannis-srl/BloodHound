-- Copyright 2024 Specter Ops, Inc.
--
-- Licensed under the Apache License, Version 2.0
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--     http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.
--
-- SPDX-License-Identifier: Apache-2.0

UPDATE asset_groups
SET tag = REGEXP_REPLACE(tag, '\s', '', 'g');

ALTER TABLE ingest_tasks
ADD COLUMN IF NOT EXISTS file_type integer DEFAULT 0;

UPDATE feature_flags SET enabled = true, user_updatable = false WHERE key = 'clear_graph_data';
