# coding=utf-8
# *** WARNING: this file was generated by pulumi. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

from . import _utilities
import typing
# Export this package's modules as members:
from .provider import *

# Make subpackages available:
if typing.TYPE_CHECKING:
    import pulumi_zeet.model as __model
    model = __model
    import pulumi_zeet.resources as __resources
    resources = __resources
    import pulumi_zeet.time as __time
    time = __time
else:
    model = _utilities.lazy_import('pulumi_zeet.model')
    resources = _utilities.lazy_import('pulumi_zeet.resources')
    time = _utilities.lazy_import('pulumi_zeet.time')

_utilities.register(
    resource_modules="""
[
 {
  "pkg": "zeet",
  "mod": "resources",
  "fqn": "pulumi_zeet.resources",
  "classes": {
   "zeet:resources:App": "App",
   "zeet:resources:Environment": "Environment",
   "zeet:resources:Project": "Project"
  }
 }
]
""",
    resource_packages="""
[
 {
  "pkg": "zeet",
  "token": "pulumi:providers:zeet",
  "fqn": "pulumi_zeet",
  "class": "Provider"
 }
]
"""
)
