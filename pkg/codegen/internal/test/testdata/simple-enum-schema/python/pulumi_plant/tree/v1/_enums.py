# coding=utf-8
# *** WARNING: this file was generated by test. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

from enum import Enum

__all__ = [
    'Diameter',
    'Farm',
    'RubberTreeVariety',
    'TreeSize',
]


class Diameter(float, Enum):
    SIXINCH = 6
    TWELVEINCH = 12


class Farm(str, Enum):
    PULUMI_PLANTERS_INC_ = "Pulumi Planters Inc."
    PLANTS_R_US = "Plants'R'Us"


class RubberTreeVariety(str, Enum):
    """
    types of rubber trees
    """
    BURGUNDY = "Burgundy"
    RUBY = "Ruby"
    TINEKE = "Tineke"


class TreeSize(str, Enum):
    SMALL = "small"
    MEDIUM = "medium"
    LARGE = "large"
