
# Adequate exception when token is expired
# Automatic Token prolongation
import logging
from typing import List, Union
from odahuflow.sdk.clients import route, model


logger = logging.getLogger(__name__)


class Model:
    """
    Can be used directly or retrieved from Models
    """

    def __init__(self, model_name):
        logger.debug(f'Start model initialization model_name={model_name}')

    def invoke(self, data: Union[List[List]]):
        """
        Invoke model
        :param data:
        :return:
        """
        raise NotImplementedError


class Models:

    def list(self) -> List:
        pass

    def init(self, model_name: str):
        pass