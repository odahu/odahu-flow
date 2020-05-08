import abc
from typing import Any, List


class Plugin(abc.ABC):

    @abc.abstractmethod
    def predict(self, *args, **kwargs) -> Any:
        """Decorator for model predict function"""
        pass


class PluginDiscovery(abc.ABC):

    @property
    @abc.abstractmethod
    def plugins(self) -> List[Plugin]:
        pass
