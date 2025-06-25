# Purpose

Package of Goakt actors which take invoke designated serverless functions on cloud providers and message back with results as JSON (and optionally logs, as a plain string). In order to consume the JSON results, make sure to use json struct tags when defining the struct to be unmarhsaled to.

## Currently implemented

    1. AWS Lambda Actor
