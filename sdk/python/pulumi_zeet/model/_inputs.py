# coding=utf-8
# *** WARNING: this file was generated by pulumi. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import copy
import warnings
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
from .. import _utilities

__all__ = [
    'CreateAppBuildInputArgs',
    'CreateAppDeployInputArgs',
    'CreateAppDockerInputArgs',
    'CreateAppEnvironmentVariableInputArgs',
    'CreateAppGithubInputArgs',
    'CreateAppResourcesInputArgs',
]

@pulumi.input_type
class CreateAppBuildInputArgs:
    def __init__(__self__, *,
                 type: pulumi.Input[str],
                 dockerfile_path: Optional[pulumi.Input[str]] = None):
        pulumi.set(__self__, "type", type)
        if dockerfile_path is not None:
            pulumi.set(__self__, "dockerfile_path", dockerfile_path)

    @property
    @pulumi.getter
    def type(self) -> pulumi.Input[str]:
        return pulumi.get(self, "type")

    @type.setter
    def type(self, value: pulumi.Input[str]):
        pulumi.set(self, "type", value)

    @property
    @pulumi.getter(name="dockerfilePath")
    def dockerfile_path(self) -> Optional[pulumi.Input[str]]:
        return pulumi.get(self, "dockerfile_path")

    @dockerfile_path.setter
    def dockerfile_path(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "dockerfile_path", value)


@pulumi.input_type
class CreateAppDeployInputArgs:
    def __init__(__self__, *,
                 deploy_target: pulumi.Input[str],
                 cluster_id: Optional[pulumi.Input[str]] = None):
        pulumi.set(__self__, "deploy_target", deploy_target)
        if cluster_id is not None:
            pulumi.set(__self__, "cluster_id", cluster_id)

    @property
    @pulumi.getter(name="deployTarget")
    def deploy_target(self) -> pulumi.Input[str]:
        return pulumi.get(self, "deploy_target")

    @deploy_target.setter
    def deploy_target(self, value: pulumi.Input[str]):
        pulumi.set(self, "deploy_target", value)

    @property
    @pulumi.getter(name="clusterId")
    def cluster_id(self) -> Optional[pulumi.Input[str]]:
        return pulumi.get(self, "cluster_id")

    @cluster_id.setter
    def cluster_id(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "cluster_id", value)


@pulumi.input_type
class CreateAppDockerInputArgs:
    def __init__(__self__, *,
                 docker_image: pulumi.Input[str]):
        pulumi.set(__self__, "docker_image", docker_image)

    @property
    @pulumi.getter(name="dockerImage")
    def docker_image(self) -> pulumi.Input[str]:
        return pulumi.get(self, "docker_image")

    @docker_image.setter
    def docker_image(self, value: pulumi.Input[str]):
        pulumi.set(self, "docker_image", value)


@pulumi.input_type
class CreateAppEnvironmentVariableInputArgs:
    def __init__(__self__, *,
                 name: pulumi.Input[str],
                 value: pulumi.Input[str],
                 sealed: Optional[pulumi.Input[bool]] = None):
        pulumi.set(__self__, "name", name)
        pulumi.set(__self__, "value", value)
        if sealed is not None:
            pulumi.set(__self__, "sealed", sealed)

    @property
    @pulumi.getter
    def name(self) -> pulumi.Input[str]:
        return pulumi.get(self, "name")

    @name.setter
    def name(self, value: pulumi.Input[str]):
        pulumi.set(self, "name", value)

    @property
    @pulumi.getter
    def value(self) -> pulumi.Input[str]:
        return pulumi.get(self, "value")

    @value.setter
    def value(self, value: pulumi.Input[str]):
        pulumi.set(self, "value", value)

    @property
    @pulumi.getter
    def sealed(self) -> Optional[pulumi.Input[bool]]:
        return pulumi.get(self, "sealed")

    @sealed.setter
    def sealed(self, value: Optional[pulumi.Input[bool]]):
        pulumi.set(self, "sealed", value)


@pulumi.input_type
class CreateAppGithubInputArgs:
    def __init__(__self__, *,
                 production_branch: pulumi.Input[str],
                 url: pulumi.Input[str]):
        pulumi.set(__self__, "production_branch", production_branch)
        pulumi.set(__self__, "url", url)

    @property
    @pulumi.getter(name="productionBranch")
    def production_branch(self) -> pulumi.Input[str]:
        return pulumi.get(self, "production_branch")

    @production_branch.setter
    def production_branch(self, value: pulumi.Input[str]):
        pulumi.set(self, "production_branch", value)

    @property
    @pulumi.getter
    def url(self) -> pulumi.Input[str]:
        return pulumi.get(self, "url")

    @url.setter
    def url(self, value: pulumi.Input[str]):
        pulumi.set(self, "url", value)


@pulumi.input_type
class CreateAppResourcesInputArgs:
    def __init__(__self__, *,
                 cpu: pulumi.Input[float],
                 memory: pulumi.Input[str],
                 ephemeral_storage: Optional[pulumi.Input[float]] = None,
                 spot_instance: Optional[pulumi.Input[bool]] = None):
        pulumi.set(__self__, "cpu", cpu)
        pulumi.set(__self__, "memory", memory)
        if ephemeral_storage is not None:
            pulumi.set(__self__, "ephemeral_storage", ephemeral_storage)
        if spot_instance is not None:
            pulumi.set(__self__, "spot_instance", spot_instance)

    @property
    @pulumi.getter
    def cpu(self) -> pulumi.Input[float]:
        return pulumi.get(self, "cpu")

    @cpu.setter
    def cpu(self, value: pulumi.Input[float]):
        pulumi.set(self, "cpu", value)

    @property
    @pulumi.getter
    def memory(self) -> pulumi.Input[str]:
        return pulumi.get(self, "memory")

    @memory.setter
    def memory(self, value: pulumi.Input[str]):
        pulumi.set(self, "memory", value)

    @property
    @pulumi.getter(name="ephemeralStorage")
    def ephemeral_storage(self) -> Optional[pulumi.Input[float]]:
        return pulumi.get(self, "ephemeral_storage")

    @ephemeral_storage.setter
    def ephemeral_storage(self, value: Optional[pulumi.Input[float]]):
        pulumi.set(self, "ephemeral_storage", value)

    @property
    @pulumi.getter(name="spotInstance")
    def spot_instance(self) -> Optional[pulumi.Input[bool]]:
        return pulumi.get(self, "spot_instance")

    @spot_instance.setter
    def spot_instance(self, value: Optional[pulumi.Input[bool]]):
        pulumi.set(self, "spot_instance", value)


