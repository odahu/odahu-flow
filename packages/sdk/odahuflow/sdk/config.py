#
#    Copyright 2019 EPAM Systems
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
"""
odahuflow env names
"""
import configparser
import logging
import os
from pathlib import Path

# Get list of all variables

ALL_VARIABLES = {}

_LOGGER = logging.getLogger()

_INI_FILE_TRIED_TO_BE_LOADED = False
_INI_FILE_CONTENT: configparser.ConfigParser = None
_INI_FILE_DEFAULT_CONFIG_PATH = Path.home().joinpath('.odahuflow/config')
_DEFAULT_INI_SECTION = 'general'


def reset_context():
    """
    Reset configuration context

    :return: None
    """
    global _INI_FILE_TRIED_TO_BE_LOADED
    global _INI_FILE_CONTENT

    _INI_FILE_TRIED_TO_BE_LOADED = False
    _INI_FILE_CONTENT = None


def get_config_file_path():
    """
    Return the config path.
    ODAHUFLOW_CONFIG environment can override path value

    :return: Path -- config path
    """
    config_path_from_env = os.getenv('ODAHUFLOW_CONFIG')

    return Path(config_path_from_env) if config_path_from_env else _INI_FILE_DEFAULT_CONFIG_PATH


def _load_config_file():
    """
    Load configuration file if it has not been loaded. Update _INI_FILE_TRIED_TO_BE_LOADED, _INI_FILE_CONTENT

    :return: None
    """
    global _INI_FILE_TRIED_TO_BE_LOADED

    if _INI_FILE_TRIED_TO_BE_LOADED:
        return

    config_path = get_config_file_path()
    _INI_FILE_TRIED_TO_BE_LOADED = True

    _LOGGER.debug('Trying to load configuration file {}'.format(config_path))

    try:
        if config_path.exists():
            config = configparser.ConfigParser()
            config.read(str(config_path))

            global _INI_FILE_CONTENT
            _INI_FILE_CONTENT = config

            _LOGGER.debug('Configuration file {} has been loaded'.format(config_path))
        else:
            _LOGGER.debug('Cannot find configuration file {}'.format(config_path))
    except Exception as exc:
        _LOGGER.exception('Cannot read config file {}'.format(config_path), exc_info=exc)


def get_config_file_section(section=_DEFAULT_INI_SECTION, silent=False):
    """
    Get section from config file

    :param section: (Optional) name of section
    :type section: str
    :param silent: (Optional) ignore if there is no file
    :type silent: bool
    :return: dict[str, str] -- values from section
    """
    _load_config_file()
    if not _INI_FILE_CONTENT:
        if silent:
            return dict()
        else:
            raise Exception('Configuration file cannot be loaded')

    if not _INI_FILE_CONTENT.has_section(section):
        return {}

    return dict(_INI_FILE_CONTENT[section])


def get_config_file_variable(variable, section=_DEFAULT_INI_SECTION):
    """
    Get variable by name from specific (or default) section

    :param variable: Name of variable
    :type variable: str
    :param section: (Optional) name of section
    :type section: str
    :return: str or None -- value
    """
    if not variable:
        return None

    _load_config_file()
    if not _INI_FILE_CONTENT:
        return None

    return _INI_FILE_CONTENT.get(section, variable, fallback=None)


def update_config_file(section=_DEFAULT_INI_SECTION, **new_values):
    """
    Update config file with new values

    :param section: (Optional) name of section to update
    :type section: str
    :param new_values: new values
    :type new_values: dict[str, typing.Optional[str]]
    :return: None
    """
    global _INI_FILE_TRIED_TO_BE_LOADED
    global _INI_FILE_CONTENT

    _load_config_file()
    config_path = get_config_file_path()

    content = _INI_FILE_CONTENT if _INI_FILE_CONTENT else configparser.ConfigParser()

    config_path.parent.mkdir(mode=0o775, parents=True, exist_ok=True)
    config_path.touch(mode=0o600, exist_ok=True)

    if not content.has_section(section):
        content.add_section(section)

    for key, value in new_values.items():
        if value:
            content.set(section, key, value)
        else:
            if section in content and key in content[section]:
                del content[section][key]

    with config_path.open('w') as config_file:
        content.write(config_file)

    _INI_FILE_TRIED_TO_BE_LOADED = True
    _INI_FILE_CONTENT = content

    reinitialize_variables()

    _LOGGER.debug('Configuration file {} has been updated'.format(config_path))


def _load_variable(name, cast_type=None, configurable_manually=True):
    """
    Load variable from config file, env. Cast it to desired type.

    :param name: name of variable
    :type name: str
    :param cast_type: (Optional) function to cast
    :type cast_type: Callable[[str], any]
    :param configurable_manually: (Optional) could this variable be configured manually or not
    :type configurable_manually: bool
    :return: Any -- variable value
    """
    value = None

    # 1st level - configuration file
    if configurable_manually:
        conf_value = get_config_file_variable(name)
        if conf_value:
            value = conf_value

    # 2nd level - env. variable
    env_value = os.environ.get(name)
    if env_value:
        value = env_value

    return cast_type(value) if value is not None else None


class ConfigVariableInformation:
    """
    Object holds information about variable (name, default value, casting function, description and etc.)
    """

    def __init__(self, name, default, cast_func, description, configurable_manually):
        """
        Build information about variable

        :param name: name of variable
        :type name: str
        :param default: default value
        :type default: Any
        :param cast_func: cast function
        :type cast_func: Callable[[str], any]
        :param description: description
        :type description: str
        :param configurable_manually: is configurable manually
        :type configurable_manually: bool
        """
        self._name = name
        self._default = default
        self._cast_func = cast_func
        self._description = description
        self._configurable_manually = configurable_manually

    @property
    def name(self):
        """
        Get name of variable

        :return: str -- name
        """
        return self._name

    @property
    def default(self):
        """
        Get default variable value

        :return: Any -- default value
        """
        return self._default

    @property
    def cast_func(self):
        """
        Get cast function (from string to desired type)

        :return: Callable[[str], any] -- casting function
        """
        return self._cast_func

    @property
    def description(self):
        """
        Get human-readable description

        :return: str -- description
        """
        return self._description

    @property
    def configurable_manually(self):
        """
        Is this variable human-configurabe?

        :return: bool -- is human configurable
        """
        return self._configurable_manually


def cast_bool(value):
    """
    Convert string to bool

    :param value: string or bool
    :type value: str or bool
    :return: bool
    """
    if value is None:
        return None

    if isinstance(value, bool):
        return value

    return value.lower() in ['true', '1', 't', 'y', 'yes']


def reinitialize_variables():
    """
    Reinitialize variables due to new ENV variables

    :return: None
    """
    for value_information in ALL_VARIABLES.values():
        explicit_value = _load_variable(value_information.name,
                                        value_information.cast_func,
                                        value_information.configurable_manually)
        value = explicit_value if explicit_value is not None else value_information.default

        globals()[value_information.name] = value


class ConfigVariableDeclaration:
    """
    Class that builds declaration of variable (and returns it's value as an instance)
    """

    def __new__(cls, name, default=None, cast_func=str, description=None, configurable_manually=True):
        """
        Create new instance

        :param name: name of variable
        :type name: str
        :param default: (Optional) default variable value [will not be passed to cast_func]
        :type default: Any
        :param cast_func: (Optional) cast function for variable value
        :type cast_func: Callable[[str], any]
        :param description: (Optional) human-readable variable description
        :type description: str
        :param configurable_manually: (Optional) can be modified by config file or CLI
        :type configurable_manually: bool
        :return: Any -- default or explicit value
        """
        information = ConfigVariableInformation(name, default, cast_func, description, configurable_manually)

        explicit_value = _load_variable(name, cast_func, configurable_manually)
        value = explicit_value if explicit_value is not None else default
        ALL_VARIABLES[information.name] = information
        return value


# Transport (HTTP)

RETRY_ATTEMPTS = ConfigVariableDeclaration(
    'RETRY_ATTEMPTS', 3, int,
    'How many retries HTTP client should make in case of transient error', True
)

BACKOFF_FACTOR = ConfigVariableDeclaration(
    'BACKOFF_FACTOR', 1, int,
    'Backoff factor for retries (See https://urllib3.readthedocs.io/en/latest/reference/urllib3.util.html)', True
)

# Verbose tracing
DEBUG = ConfigVariableDeclaration('DEBUG', False, cast_bool,
                                  'Enable verbose program output',
                                  True)

# Model invocation testing
MODEL_SERVER_URL = ConfigVariableDeclaration('MODEL_SERVER_URL', '', str, 'Default url of model server', True)

MODEL_HOST = ConfigVariableDeclaration('MODEL_HOST', '', str, 'Default host of model server', True)

MODEL_DEPLOYMENT_NAME = ConfigVariableDeclaration('MODEL_DEPLOYMENT_NAME', '', str, 'Model deployment name', True)

MODEL_ROUTE_NAME = ConfigVariableDeclaration('MODEL_ROUTE_NAME', '', str, 'Model route name', True)

# API endpoint
API_URL = ConfigVariableDeclaration('API_URL', 'http://localhost:5000', str,
                                    'URL of API server',
                                    True)
API_TOKEN = ConfigVariableDeclaration('API_TOKEN', None, str,
                                      'Token for API server authorisation',
                                      True)
API_REFRESH_TOKEN = ConfigVariableDeclaration('API_REFRESH_TOKEN', None, str,
                                              'Refresh token',
                                              True)
API_ISSUING_URL = ConfigVariableDeclaration('API_ISSUING_URL', None, str,
                                            'URL for refreshing and issuing tokens',
                                            True)

ISSUER_URL = ConfigVariableDeclaration('ISSUER_URL', None, str, 'OIDC Issuer URL', True)

# Auth
ODAHUFLOWCTL_OAUTH_CLIENT_ID = ConfigVariableDeclaration('ODAHUFLOWCTL_OAUTH_CLIENT_ID', 'legion-cli', str,
                                                         'Set OAuth2 Client id',
                                                         True)

ODAHUFLOWCTL_OAUTH_CLIENT_SECRET = ConfigVariableDeclaration('ODAHUFLOWCTL_OAUTH_CLIENT_SECRET', '', str,
                                                             'Set OAuth2 Client secret',
                                                             True)

ODAHUFLOWCTL_OAUTH_SCOPE = ConfigVariableDeclaration('ODAHUFLOWCTL_OAUTH_SCOPE',
                                                     'openid profile email offline_access groups', str,
                                                     'Set OAuth2 scope',
                                                     True)

ODAHUFLOWCTL_OAUTH_LOOPBACK_HOST = ConfigVariableDeclaration('ODAHUFLOWCTL_OAUTH_LOOPBACK_HOST',
                                                             '127.0.0.1', str,
                                                             'Target redirect for OAuth2 interactive authorization',
                                                             True)

ODAHUFLOWCTL_OAUTH_LOOPBACK_URL = ConfigVariableDeclaration('ODAHUFLOWCTL_OAUTH_LOOPBACK_URL',
                                                            '/oauth/callback', str,
                                                            'Target redirect url for OAuth2 interactive authorization',
                                                            True)

ODAHUFLOWCTL_OAUTH_TOKEN_ISSUING_URL = ConfigVariableDeclaration('ODAHUFLOWCTL_OAUTH_TOKEN_ISSUING_URL',
                                                                 '', str,
                                                                 'OAuth2 token issuing URL',
                                                                 True)

ODAHUFLOWCTL_OAUTH_AUTH_URL = ConfigVariableDeclaration('ODAHUFLOWCTL_OAUTH_AUTH_URL',
                                                        '',
                                                        str,
                                                        'OAuth2 authorization URL',
                                                        True)

JUPYTER_REDIRECT_URL = ConfigVariableDeclaration('JUPYTER_REDIRECT_URL',
                                                 '', str,
                                                 'JupyterLab external URL',
                                                 True)

ODAHUFLOWCTL_NONINTERACTIVE = ConfigVariableDeclaration('ODAHUFLOWCTL_NONINTERACTIVE', False,
                                                        bool, 'Disable any interaction (e.g. authorization)', True)

# Local

LOCAL_MODEL_OUTPUT_DIR = ConfigVariableDeclaration('LOCAL_MODEL_OUTPUT_DIR',
                                                   str(Path.home().joinpath(".odahuflow", "training_output")),
                                                   str, 'Directory where model artifacts will be saved', True)
