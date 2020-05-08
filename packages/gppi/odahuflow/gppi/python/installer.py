import os

from odahuflow.gppi.model.installer import ModelInstaller


class PythonModelInstaller(ModelInstaller):

    def install(self):
        pass


def kek():
    print(os.path.join(os.path.dirname(__file__), 'library_template'))
