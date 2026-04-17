#!/usr/bin/env python3

import sys

import yaml

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("usage: {sys.argv[0]} EMAIL-FILE OUTPUT-FILE")
        sys.exit(1)

    source_file = sys.argv[1]
    group_name = source_file.split(".")[0] if "." in source_file else source_file
    entries = {"name": group_name, "users": []}
    with open(sys.argv[1], "r") as f:
        lines = f.read().split("\n")
        lines = filter(lambda l: l != "", map(lambda l: l.strip(), lines))
        for email in sorted(lines):
            email = email.lower()
            (username, _) = email.split("@")
            entry = {"name": username, "email": email, "ssh-key": ""}
            entries["users"].append(entry)

    with open(sys.argv[2], "w") as f:
        f.write(yaml.dump(entries, allow_unicode=True))
