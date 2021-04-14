#  Copyright 2021 EPAM Systems
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#  http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
from odahuflow.sdk.models import ModelTraining, ToolchainIntegration, ModelDeployment, ModelRoute, Connection, \
    ModelPackaging, PackagingIntegration, InferenceService, InferenceJob

ROOT_MODELS = {
    ModelTraining.__name__: ModelTraining,
    ToolchainIntegration.__name__: ToolchainIntegration,
    ModelDeployment.__name__: ModelDeployment,
    ModelRoute.__name__: ModelRoute,
    Connection.__name__: Connection,
    ModelPackaging.__name__: ModelPackaging,
    PackagingIntegration.__name__: PackagingIntegration,
    InferenceService.__name__: InferenceService,
    InferenceJob.__name__: InferenceJob,
}
