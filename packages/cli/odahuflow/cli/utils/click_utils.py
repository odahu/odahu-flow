#  Copyright 2020 EPAM Systems
#
#  Licensed under the Apache License, Ver
#  you may not use this file except in co
#  You may obtain a copy of the License a
#
#      http://www.apache.org/licenses/LIC
#
#  Unless required by applicable law or a
#  distributed under the License is distr
#  WITHOUT WARRANTIES OR CONDITIONS OF AN
#  See the License for the specific langu
#  limitations under the License.
#
#  Licensed under the Apache License, Ver
#  you may not use this file except in co
#  You may obtain a copy of the License a
#
#      http://www.apache.org/licenses/LIC
#
#  Unless required by applicable law or a
#  distributed under the License is distr
#  WITHOUT WARRANTIES OR CONDITIONS OF AN
#  See the License for the specific langu
#  limitations under the License.
#
#  Licensed under the Apache License, Ver
#  you may not use this file except in co
#  You may obtain a copy of the License a
#
#      http://www.apache.org/licenses/LIC
#
#  Unless required by applicable law or a
#  distributed under the License is distr
#  WITHOUT WARRANTIES OR CONDITIONS OF AN
#  See the License for the specific langu
#  limitations under the License.

import click


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
