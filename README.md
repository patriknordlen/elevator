# elevator

A utility for temporarily escalating users' privileges in a Google Cloud environment.

It works similarly to [jit-access](https://github.com/GoogleCloudPlatform/jit-access) but

- it is built in Go
- uses (currently) a local config file for defining escalation policies instead of IAM conditions
  - this means it's quicker in determining applicable policies, but doesn't support policy inheritance
- aims to also support escalations on the folder and org level

Currently in an early PoC stage.