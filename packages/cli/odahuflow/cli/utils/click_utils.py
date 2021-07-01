#  Copyright 2020 EPAM Systems
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
from functools import wraps

import click

from odahuflow.sdk import config
from odahuflow.sdk.clients.api import RemoteAPIClient

ABBREVIATION = {
    "conn": "connection",
    "conf": "config",
    "dep": "deployment",
    "pack": "packaging",
    "pi": "packaging-integration",
    "temp": "template",
    "ti": "toolchain-integration",
    "train": "training"
}


class BetterHelpGroup(click.Group):
    """
    This class allows user to get a subcommand's help text without invocation of group's callback

    This is to avoid bad user experience when user wants to check a subcommand's help text, but CLI asks him to
    fulfill all required arguments of a group. Also even with all required arguments in place, group's callback
    can fail for any reason, which is non-sense when asking CLI for a subcommand's help.
    """

    def invoke(self, ctx: click.Context):
        """
        This method sets a callback function (the one that decorated by click.group()) to None,
        if --help appears in args. So the callback is ignored at all if user asks to print help text.
        """
        if any(arg in ctx.help_option_names for arg in ctx.args):
            self.callback = None
        return super().invoke(ctx)


class AbbreviationGroup(BetterHelpGroup):
    """
    AbbreviationGroup
    """

    def get_command(self, ctx, cmd_name):
        """
        Override get command of click.Group
        :param ctx: click context
        :param cmd_name: group name
        :return: click Command
        """
        rv = click.Group.get_command(self, ctx, cmd_name)
        if rv is not None:
            return rv

        return click.Group.get_command(self, ctx, ABBREVIATION.get(cmd_name))


def auth_options(handler_f):

    @click.option('--url', help='API server host', default=config.API_URL)
    @click.option('--client-id', help='Client ID for authorization', default=config.ODAHUFLOWCTL_OAUTH_CLIENT_ID)
    @click.option('--client-secret', help='Client Secret for authorization',
                  default=config.ODAHUFLOWCTL_OAUTH_CLIENT_SECRET)
    @click.option('--issuer-url', help='Token Issuer URL (Identity Provider)', default=config.ISSUER_URL)
    @click.option('--token', help='API server jwt token', default=config.API_TOKEN)
    @wraps(handler_f)
    def wrapper(*args, url, client_id, client_secret, issuer_url, token, **kwargs):
        base_client = RemoteAPIClient(
            base_url=url,
            client_id=client_id,
            client_secret=client_secret,
            issuer_url=issuer_url,
            token=token
        )

        handler_f(*args, api_client=base_client, **kwargs)

    return wrapper
