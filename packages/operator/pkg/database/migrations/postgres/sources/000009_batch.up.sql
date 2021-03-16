/*
 *
 *     Copyright 2021 EPAM Systems
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 */

BEGIN;

CREATE TABLE IF NOT EXISTS odahu_batch_inference_service
(
    id   VARCHAR(64) PRIMARY KEY,
    created timestamptz,
    updated timestamptz,
    deletionmark boolean default FALSE not null,
    spec JSONB,
    status JSONB
);

CREATE TABLE IF NOT EXISTS odahu_batch_inference_job
(
    id   VARCHAR(64) PRIMARY KEY,
    created timestamptz not null,
    updated timestamptz not null,
    deletionmark boolean default FALSE not null,
    spec JSONB not null,
    status JSONB not null,
    service      varchar(64)           not null
    constraint odahu_bij_bis_fk
    references odahu_batch_inference_service
    on update restrict on delete restrict
);

COMMIT;