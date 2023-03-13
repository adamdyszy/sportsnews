# Welcome to contributing to sports news!

Here you can get to know how the project is structured and how to add features.

Try to follow twelve-factor methodology https://12factor.net/ for good practices. 

## Project structure

Source files:

- api - hosts all HTTP related files, it is a good practice to group and version the code. Under this you’ll find the router and handler for the server logic.
- cmd/<project-name> - keeps the main.go, keep this file relatively small as a good practice.
  - config - next to the main.go include config with default configuration
- internal - folder will have all the different logic packages such as storage implementation
- storage - overall storage interfaces
- types - api used types

Test files will be included in packages they correspond to.

## Dependencies

Use go mod vendor to force reproducibility. For future using proxy with replaces might be a good option https://proxy.golang.org/

## Config

Use viper for configuration

## Documentation

Add comments above function, constants, packages to document it.
To run docs locally install godoc and run it on localhost:8080 using:

```bash
go install golang.org/x/tools/cmd/godoc@latest
godoc -http :8080
```