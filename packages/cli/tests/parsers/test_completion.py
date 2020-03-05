#
#    Copyright 2020 EPAM Systems
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
#
import pytest
from click.testing import CliRunner
from click_completion import Shell
from odahuflow.cli.parsers import completion


@pytest.mark.parametrize("shell", [e.name for e in Shell])
def test_delete_by_file(shell: str):
    result = CliRunner().invoke(completion.completion, [shell])

    assert result.exit_code == 0
    assert result.stdout
