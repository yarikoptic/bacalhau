# coding: utf-8

"""
    Bacalhau API

    This page is the reference of the Bacalhau REST API. Project docs are available at https://docs.bacalhau.org/. Find more information about Bacalhau at https://github.com/filecoin-project/bacalhau.  # noqa: E501

    OpenAPI spec version: 0.3.22.post4
    Contact: team@bacalhau.org
    Generated by: https://github.com/swagger-api/swagger-codegen.git
"""


import pprint
import re  # noqa: F401

import six

from bacalhau_apiclient.configuration import Configuration


class ShardState(object):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    """
    Attributes:
      swagger_types (dict): The key is attribute name
                            and the value is attribute type.
      attribute_map (dict): The key is attribute name
                            and the value is json key in definition.
    """
    swagger_types = {
        'create_time': 'str',
        'executions': 'list[ExecutionState]',
        'job_id': 'str',
        'shard_index': 'int',
        'state': 'ShardStateState',
        'update_time': 'str',
        'version': 'int'
    }

    attribute_map = {
        'create_time': 'CreateTime',
        'executions': 'Executions',
        'job_id': 'JobID',
        'shard_index': 'ShardIndex',
        'state': 'State',
        'update_time': 'UpdateTime',
        'version': 'Version'
    }

    def __init__(self, create_time=None, executions=None, job_id=None, shard_index=None, state=None, update_time=None, version=None, _configuration=None):  # noqa: E501
        """ShardState - a model defined in Swagger"""  # noqa: E501
        if _configuration is None:
            _configuration = Configuration()
        self._configuration = _configuration

        self._create_time = None
        self._executions = None
        self._job_id = None
        self._shard_index = None
        self._state = None
        self._update_time = None
        self._version = None
        self.discriminator = None

        if create_time is not None:
            self.create_time = create_time
        if executions is not None:
            self.executions = executions
        if job_id is not None:
            self.job_id = job_id
        if shard_index is not None:
            self.shard_index = shard_index
        if state is not None:
            self.state = state
        if update_time is not None:
            self.update_time = update_time
        if version is not None:
            self.version = version

    @property
    def create_time(self):
        """Gets the create_time of this ShardState.  # noqa: E501

        CreateTime is the time when the shard was created, which is the same as the job creation time.  # noqa: E501

        :return: The create_time of this ShardState.  # noqa: E501
        :rtype: str
        """
        return self._create_time

    @create_time.setter
    def create_time(self, create_time):
        """Sets the create_time of this ShardState.

        CreateTime is the time when the shard was created, which is the same as the job creation time.  # noqa: E501

        :param create_time: The create_time of this ShardState.  # noqa: E501
        :type: str
        """

        self._create_time = create_time

    @property
    def executions(self):
        """Gets the executions of this ShardState.  # noqa: E501

        Executions is a list of executions of the shard across the nodes. A new execution is created when a node is selected to execute the shard, and a node can have multiple executions for the same shard due to retries, but there can only be a single active execution per node at any given time.  # noqa: E501

        :return: The executions of this ShardState.  # noqa: E501
        :rtype: list[ExecutionState]
        """
        return self._executions

    @executions.setter
    def executions(self, executions):
        """Sets the executions of this ShardState.

        Executions is a list of executions of the shard across the nodes. A new execution is created when a node is selected to execute the shard, and a node can have multiple executions for the same shard due to retries, but there can only be a single active execution per node at any given time.  # noqa: E501

        :param executions: The executions of this ShardState.  # noqa: E501
        :type: list[ExecutionState]
        """

        self._executions = executions

    @property
    def job_id(self):
        """Gets the job_id of this ShardState.  # noqa: E501

        JobID is the unique identifier for the job  # noqa: E501

        :return: The job_id of this ShardState.  # noqa: E501
        :rtype: str
        """
        return self._job_id

    @job_id.setter
    def job_id(self, job_id):
        """Sets the job_id of this ShardState.

        JobID is the unique identifier for the job  # noqa: E501

        :param job_id: The job_id of this ShardState.  # noqa: E501
        :type: str
        """

        self._job_id = job_id

    @property
    def shard_index(self):
        """Gets the shard_index of this ShardState.  # noqa: E501

        ShardIndex is the index of the shard in the job  # noqa: E501

        :return: The shard_index of this ShardState.  # noqa: E501
        :rtype: int
        """
        return self._shard_index

    @shard_index.setter
    def shard_index(self, shard_index):
        """Sets the shard_index of this ShardState.

        ShardIndex is the index of the shard in the job  # noqa: E501

        :param shard_index: The shard_index of this ShardState.  # noqa: E501
        :type: int
        """

        self._shard_index = shard_index

    @property
    def state(self):
        """Gets the state of this ShardState.  # noqa: E501


        :return: The state of this ShardState.  # noqa: E501
        :rtype: ShardStateState
        """
        return self._state

    @state.setter
    def state(self, state):
        """Sets the state of this ShardState.


        :param state: The state of this ShardState.  # noqa: E501
        :type: ShardStateState
        """

        self._state = state

    @property
    def update_time(self):
        """Gets the update_time of this ShardState.  # noqa: E501

        UpdateTime is the time when the shard state was last updated.  # noqa: E501

        :return: The update_time of this ShardState.  # noqa: E501
        :rtype: str
        """
        return self._update_time

    @update_time.setter
    def update_time(self, update_time):
        """Sets the update_time of this ShardState.

        UpdateTime is the time when the shard state was last updated.  # noqa: E501

        :param update_time: The update_time of this ShardState.  # noqa: E501
        :type: str
        """

        self._update_time = update_time

    @property
    def version(self):
        """Gets the version of this ShardState.  # noqa: E501

        Version is the version of the shard state. It is incremented every time the shard state is updated.  # noqa: E501

        :return: The version of this ShardState.  # noqa: E501
        :rtype: int
        """
        return self._version

    @version.setter
    def version(self, version):
        """Sets the version of this ShardState.

        Version is the version of the shard state. It is incremented every time the shard state is updated.  # noqa: E501

        :param version: The version of this ShardState.  # noqa: E501
        :type: int
        """

        self._version = version

    def to_dict(self):
        """Returns the model properties as a dict"""
        result = {}

        for attr, _ in six.iteritems(self.swagger_types):
            value = getattr(self, attr)
            if isinstance(value, list):
                result[attr] = list(map(
                    lambda x: x.to_dict() if hasattr(x, "to_dict") else x,
                    value
                ))
            elif hasattr(value, "to_dict"):
                result[attr] = value.to_dict()
            elif isinstance(value, dict):
                result[attr] = dict(map(
                    lambda item: (item[0], item[1].to_dict())
                    if hasattr(item[1], "to_dict") else item,
                    value.items()
                ))
            else:
                result[attr] = value
        if issubclass(ShardState, dict):
            for key, value in self.items():
                result[key] = value

        return result

    def to_str(self):
        """Returns the string representation of the model"""
        return pprint.pformat(self.to_dict())

    def __repr__(self):
        """For `print` and `pprint`"""
        return self.to_str()

    def __eq__(self, other):
        """Returns true if both objects are equal"""
        if not isinstance(other, ShardState):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, ShardState):
            return True

        return self.to_dict() != other.to_dict()
