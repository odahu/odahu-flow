#
#    Copyright 2017 EPAM Systems
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
from setuptools import setup
import os
import re

PACKAGE_ROOT_PATH = os.path.dirname(os.path.abspath(__file__))


def extract_requirements(filename):
    """
    Extracts requirements from a pip formatted requirements file.
    :param filename: str path to file
    :return: list of package names as strings
    """
    with open(filename, 'r') as requirements_file:
        return requirements_file.read().splitlines()


def extract_version(filename):
    """
    Extract version from .py file using regex
    :param filename: str path to file
    :return: str version
    """
    with open(filename, 'rt') as version_file:
        file_content = version_file.read()
        VSRE = r"^__version__ = ['\"]([^'\"]*)['\"]"
        mo = re.search(VSRE, file_content, re.M)
        if mo:
            return mo.group(1)
        else:
            raise RuntimeError("Unable to find version string in %s." % (file_content,))


setup(name='drun_test_tools',
      version=extract_version(os.path.join(PACKAGE_ROOT_PATH, 'drun_test_tools', 'version.py')),
      description='DRUN CI tools',
      url='http://github.com/akharlamov/drun-root',
      author='Kirill Makhonin',
      author_email='kirill_makhonin@epam.com',
      license='Apache v2',
      packages=['drun_test_tools'],
      include_package_data=True,
      scripts=['bin/create_example_jobs',
               'bin/drun_bootstrap_grafana',
               'bin/check_jenkins_jobs',
               'bin/update_version_id',
               'bin/copyright_scanner'],
      install_requires=extract_requirements(os.path.join(PACKAGE_ROOT_PATH, 'requirements', 'base.txt')),
      test_suite='nose.collector',
      tests_require=extract_requirements(os.path.join(PACKAGE_ROOT_PATH, 'requirements', 'test.txt')),
      zip_safe=False)
