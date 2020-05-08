import abc
from typing import List

from odahuflow.gppi.model.model import Model


class ModelDiscovery(abc.ABC):

    @property
    @abc.abstractmethod
    def models(self) -> List[Model]:
        pass

    @abc.abstractmethod
    def get_model(self, entrypoint_model_name: str) -> Model:
        """
        TODO: can here insert proxy classes plugins
        :param entrypoint_model_name:
        """
        pass
