# Friend Quest Ark Controller
This repository is broken up into three separate parts: Core, CLI, API.

---

## Core
The Core contains basic interaction functions used to communicate with AWS and RCON.

Environment variables can be set in a .env file at the root of the repository during testing. Please find the list of required variables below.

***Required Environment variables***
* ACSRCONPASS - RCON Password for the Ark Servers

---

## CLI
The CLI uses the functions from Core and makes them available via a CLI to make local control easier.

---

## API 
The API is an AWS Lambda mirco service paired withan AWS S3 bucket intended to provide a GUI for the Core functions.
