This program provides a set of functions for template:

> Using https://github.com/Masterminds/sprig as default template functions

- "toToml": Converts the input to TOML (Tom's Obvious, Minimal Language) format.
  Usage: toToml "your_input"

- "toYaml": Converts the input to YAML (YAML Ain't Markup Language) format.
  Usage: toYaml "your_input"

- "fromYaml": Converts a YAML format input to a standard string.
  Usage: fromYaml "your_yaml_input"

- "fromYamlArray": Converts a YAML format array to a standard array.
  Usage: fromYamlArray "your_yaml_array"

- "toJson": Converts the input to JSON format.
  Usage: toJson "your_input"

- "fromJson": Converts a JSON format input to a standard string.
  Usage: fromJson "your_json_input"

- "fromJsonArray": Converts a JSON format array to a standard array.
  Usage: fromJsonArray "your_json_array"

- "ipNet": Returns the network IP for a given IP address.
  Usage: ipNet "your_ip_address"

- "ipAt": Returns the IP address at a specific position in a network range.
  Usage: ipAt "your_network_range" "position"

- "readFile": Reads a file and returns its content.
  Usage: readFile "your_file_path"
