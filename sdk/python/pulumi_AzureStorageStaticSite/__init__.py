# coding=utf-8
# *** WARNING: this file was generated by Pulumi SDK Generator. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

from . import _utilities
import typing
# Export this package's modules as members:
from .provider import *
from .static_page import *
_utilities.register(
    resource_modules="""
[
 {
  "pkg": "AzureStorageStaticSite",
  "mod": "index",
  "fqn": "pulumi_AzureStorageStaticSite",
  "classes": {
   "AzureStorageStaticSite:index:StaticPage": "StaticPage"
  }
 }
]
""",
    resource_packages="""
[
 {
  "pkg": "AzureStorageStaticSite",
  "token": "pulumi:providers:AzureStorageStaticSite",
  "fqn": "pulumi_AzureStorageStaticSite",
  "class": "Provider"
 }
]
"""
)
