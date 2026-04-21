#!/usr/bin/env python3
"""Convert Terraform provider schema from protocol v6 nested_type format
to the protocol v5 block_types format that Upjet v1.11.0 expects.

The incident-io/incident Terraform provider uses protocol v6 nested attributes
(nested_type with nesting_mode), but Upjet v1.11.0 expects the older block_types
format. This script converts between the two formats.

Usage:
    python3 hack/convert-schema.py config/schema.json config/schema.json
"""

import json
import sys
import copy


def convert_nested_type_to_block_type(attr_schema):
    """Convert a nested_type attribute to a block_types entry."""
    nested = attr_schema["nested_type"]
    nesting_mode = nested.get("nesting_mode", "single")

    # Map nesting_mode to block_types nesting_mode
    # Note: Upjet v1.11.0 doesn't handle "map" nesting well, convert to "list"
    mode_map = {
        "single": "single",
        "list": "list",
        "set": "set",
        "map": "list",
    }

    block = {"attributes": {}, "description": attr_schema.get("description", ""),
             "description_kind": attr_schema.get("description_kind", "plain")}

    # Process nested attributes
    for name, nested_attr in nested.get("attributes", {}).items():
        if "nested_type" in nested_attr:
            # Recursively convert nested nested_types
            if "block_types" not in block:
                block["block_types"] = {}
            block["block_types"][name] = convert_nested_type_to_block_type(
                nested_attr
            )
        else:
            block["attributes"][name] = nested_attr

    result = {
        "nesting_mode": mode_map.get(nesting_mode, "list"),
        "block": block,
    }

    # Add min/max items based on nesting mode and required status
    if nesting_mode == "single":
        result["max_items"] = 1
        if attr_schema.get("required"):
            result["min_items"] = 1
    elif nesting_mode in ("list", "set"):
        if attr_schema.get("required"):
            result["min_items"] = 1

    return result


def convert_resource_schema(resource_schema):
    """Convert a resource schema, moving nested_type attrs to block_types."""
    schema = copy.deepcopy(resource_schema)
    block = schema.get("block", {})

    attrs = block.get("attributes", {})
    block_types = block.get("block_types", {})

    attrs_to_remove = []
    for name, attr in attrs.items():
        if "nested_type" in attr:
            block_types[name] = convert_nested_type_to_block_type(attr)
            attrs_to_remove.append(name)

    for name in attrs_to_remove:
        del attrs[name]

    if block_types:
        block["block_types"] = block_types

    schema["block"] = block
    return schema


def convert_schema(schema):
    """Convert the full provider schema."""
    schema = copy.deepcopy(schema)

    for provider_key, provider_schema in schema.get("provider_schemas", {}).items():
        # Convert resource schemas
        for res_name, res_schema in provider_schema.get(
            "resource_schemas", {}
        ).items():
            provider_schema["resource_schemas"][res_name] = convert_resource_schema(
                res_schema
            )

        # Convert data source schemas
        for ds_name, ds_schema in provider_schema.get(
            "data_source_schemas", {}
        ).items():
            provider_schema["data_source_schemas"][ds_name] = convert_resource_schema(
                ds_schema
            )

    return schema


def main():
    if len(sys.argv) < 3:
        print(f"Usage: {sys.argv[0]} <input.json> <output.json>")
        sys.exit(1)

    input_path = sys.argv[1]
    output_path = sys.argv[2]

    with open(input_path) as f:
        schema = json.load(f)

    converted = convert_schema(schema)

    with open(output_path, "w") as f:
        json.dump(converted, f, indent=None)

    # Count conversions
    total = 0
    for pkey, pschema in converted.get("provider_schemas", {}).items():
        for rname, rschema in pschema.get("resource_schemas", {}).items():
            bt = rschema.get("block", {}).get("block_types", {})
            total += len(bt)

    print(f"Converted schema written to {output_path}")
    print(f"Total block_types entries: {total}")


if __name__ == "__main__":
    main()
