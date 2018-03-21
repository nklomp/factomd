"""
Module for reading and validating the network configuration file.
"""
from collections import namedtuple
from schema import Schema, SchemaError, Optional, Or
import yaml

from nettool import log


NODE = Schema({
    "name": str,
    Optional("seed"): bool,
    Optional("role"): Or("follower", "federated", "audit"),
    Optional("ui_port"): int,
    Optional("api_port"): int,
    Optional("flags"): str
})


RULE = Schema({
    "action": Or("allow", "deny"),
    Optional("source"): str,
    Optional("target"): str,
    Optional("one-way"): bool
})


CONFIG = Schema({
    "nodes": [NODE],
    "network": {
        "rules": [RULE]
    },
    Optional("flags"): str
})


Environment = namedtuple("Environment", "flags, nodes, network")

Node = namedtuple("Node", "name, seed, role, ui_port, api_port, flags")

Network = namedtuple("Network", "rules")

Rule = namedtuple("Rule", "source, target, action")


def read_file(config_path):
    """
    Reads the network setup from the config file.
    """
    cfg = _read_yaml(config_path)
    _validate_schema(cfg)
    return _parse_env_config(cfg)


def _read_yaml(path):
    with open(path) as net_file:
        return yaml.load(net_file)


def _validate_schema(cfg):
    try:
        CONFIG.validate(cfg)
    except SchemaError as exc:
        log.fatal(exc)


def _parse_env_config(cfg):
    env = Environment(
        nodes=[_parse_node(node) for node in cfg["nodes"]],
        network=_parse_network(cfg["network"]),
        flags=cfg.get("flags", None)
    )
    if not env.nodes:
        log.fatal("At least one node needs to be defined in the config file")

    return env


def _parse_node(cfg):
    return Node(
        name=cfg["name"],
        seed=cfg.get("seed", False),
        role=cfg.get("role", "follower"),
        ui_port=cfg.get("ui_port", None),
        api_port=cfg.get("api_port", None),
        flags=cfg.get("flags", None))


def _parse_network(cfg):
    rules = []

    for rule_cfg in cfg["rules"]:
        source = rule_cfg.get("source", "*")
        target = rule_cfg.get("target", "*")
        action = rule_cfg.get("action", "deny")
        one_way = rule_cfg.get("one-way", False)

        rules.append(Rule(source, target, action))
        if not one_way:
            rules.append(Rule(target, source, action))

    return Network(rules=rules)
