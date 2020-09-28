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
import click
import click_completion
import click_completion.core


@click.command()
@click.argument('shell', required=False, type=click_completion.DocumentedChoice(
    click_completion.core.shells))
def completion(shell):
    """
    Output odahuflowctl completion code to stdout.\n
    \b
    Load zsh completion in the current shell:
        source <(odahuflowctl completion zsh)
    \b
    Load powershell completion in the current shell:
        odahuflowctl completion > .\odahu_completion.ps1;
        . .\odahu_completion.ps1;
        Remove-Item .\odahu_completion.ps1
    \f
    """
    click.echo(click_completion.core.get_code(shell))
