# coding: utf-8

"""
    Bacalhau API

    This page is the reference of the Bacalhau REST API. Project docs are available at https://docs.bacalhau.org/. Find more information about Bacalhau at https://github.com/bacalhau-project/bacalhau.  # noqa: E501

    OpenAPI spec version: 0.3.23.post7
    Contact: team@bacalhau.org
    Generated by: https://github.com/swagger-api/swagger-codegen.git
"""

import pprint
import re  # noqa: F401

import six


class ListRequest(object):
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
        "client_id": "str",
        "exclude_tags": "list[str]",
        "id": "str",
        "include_tags": "list[str]",
        "max_jobs": "int",
        "return_all": "bool",
        "sort_by": "str",
        "sort_reverse": "bool",
    }

    attribute_map = {
        "client_id": "client_id",
        "exclude_tags": "exclude_tags",
        "id": "id",
        "include_tags": "include_tags",
        "max_jobs": "max_jobs",
        "return_all": "return_all",
        "sort_by": "sort_by",
        "sort_reverse": "sort_reverse",
    }

    def __init__(
        self,
        client_id=None,
        exclude_tags=None,
        id=None,
        include_tags=None,
        max_jobs=None,
        return_all=None,
        sort_by=None,
        sort_reverse=None,
    ):  # noqa: E501
        """ListRequest - a model defined in Swagger"""  # noqa: E501
        self._client_id = None
        self._exclude_tags = None
        self._id = None
        self._include_tags = None
        self._max_jobs = None
        self._return_all = None
        self._sort_by = None
        self._sort_reverse = None
        self.discriminator = None
        if client_id is not None:
            self.client_id = client_id
        if exclude_tags is not None:
            self.exclude_tags = exclude_tags
        if id is not None:
            self.id = id
        if include_tags is not None:
            self.include_tags = include_tags
        if max_jobs is not None:
            self.max_jobs = max_jobs
        if return_all is not None:
            self.return_all = return_all
        if sort_by is not None:
            self.sort_by = sort_by
        if sort_reverse is not None:
            self.sort_reverse = sort_reverse

    @property
    def client_id(self):
        """Gets the client_id of this ListRequest.  # noqa: E501


        :return: The client_id of this ListRequest.  # noqa: E501
        :rtype: str
        """
        return self._client_id

    @client_id.setter
    def client_id(self, client_id):
        """Sets the client_id of this ListRequest.


        :param client_id: The client_id of this ListRequest.  # noqa: E501
        :type: str
        """

        self._client_id = client_id

    @property
    def exclude_tags(self):
        """Gets the exclude_tags of this ListRequest.  # noqa: E501


        :return: The exclude_tags of this ListRequest.  # noqa: E501
        :rtype: list[str]
        """
        return self._exclude_tags

    @exclude_tags.setter
    def exclude_tags(self, exclude_tags):
        """Sets the exclude_tags of this ListRequest.


        :param exclude_tags: The exclude_tags of this ListRequest.  # noqa: E501
        :type: list[str]
        """

        self._exclude_tags = exclude_tags

    @property
    def id(self):
        """Gets the id of this ListRequest.  # noqa: E501


        :return: The id of this ListRequest.  # noqa: E501
        :rtype: str
        """
        return self._id

    @id.setter
    def id(self, id):
        """Sets the id of this ListRequest.


        :param id: The id of this ListRequest.  # noqa: E501
        :type: str
        """

        self._id = id

    @property
    def include_tags(self):
        """Gets the include_tags of this ListRequest.  # noqa: E501


        :return: The include_tags of this ListRequest.  # noqa: E501
        :rtype: list[str]
        """
        return self._include_tags

    @include_tags.setter
    def include_tags(self, include_tags):
        """Sets the include_tags of this ListRequest.


        :param include_tags: The include_tags of this ListRequest.  # noqa: E501
        :type: list[str]
        """

        self._include_tags = include_tags

    @property
    def max_jobs(self):
        """Gets the max_jobs of this ListRequest.  # noqa: E501


        :return: The max_jobs of this ListRequest.  # noqa: E501
        :rtype: int
        """
        return self._max_jobs

    @max_jobs.setter
    def max_jobs(self, max_jobs):
        """Sets the max_jobs of this ListRequest.


        :param max_jobs: The max_jobs of this ListRequest.  # noqa: E501
        :type: int
        """

        self._max_jobs = max_jobs

    @property
    def return_all(self):
        """Gets the return_all of this ListRequest.  # noqa: E501


        :return: The return_all of this ListRequest.  # noqa: E501
        :rtype: bool
        """
        return self._return_all

    @return_all.setter
    def return_all(self, return_all):
        """Sets the return_all of this ListRequest.


        :param return_all: The return_all of this ListRequest.  # noqa: E501
        :type: bool
        """

        self._return_all = return_all

    @property
    def sort_by(self):
        """Gets the sort_by of this ListRequest.  # noqa: E501


        :return: The sort_by of this ListRequest.  # noqa: E501
        :rtype: str
        """
        return self._sort_by

    @sort_by.setter
    def sort_by(self, sort_by):
        """Sets the sort_by of this ListRequest.


        :param sort_by: The sort_by of this ListRequest.  # noqa: E501
        :type: str
        """

        self._sort_by = sort_by

    @property
    def sort_reverse(self):
        """Gets the sort_reverse of this ListRequest.  # noqa: E501


        :return: The sort_reverse of this ListRequest.  # noqa: E501
        :rtype: bool
        """
        return self._sort_reverse

    @sort_reverse.setter
    def sort_reverse(self, sort_reverse):
        """Sets the sort_reverse of this ListRequest.


        :param sort_reverse: The sort_reverse of this ListRequest.  # noqa: E501
        :type: bool
        """

        self._sort_reverse = sort_reverse

    def to_dict(self):
        """Returns the model properties as a dict"""
        result = {}

        for attr, _ in six.iteritems(self.swagger_types):
            value = getattr(self, attr)
            if isinstance(value, list):
                result[attr] = list(
                    map(lambda x: x.to_dict() if hasattr(x, "to_dict") else x, value)
                )
            elif hasattr(value, "to_dict"):
                result[attr] = value.to_dict()
            elif isinstance(value, dict):
                result[attr] = dict(
                    map(
                        lambda item: (item[0], item[1].to_dict())
                        if hasattr(item[1], "to_dict")
                        else item,
                        value.items(),
                    )
                )
            else:
                result[attr] = value
        if issubclass(ListRequest, dict):
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
        if not isinstance(other, ListRequest):
            return False

        return self.__dict__ == other.__dict__

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        return not self == other
