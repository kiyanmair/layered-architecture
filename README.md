# layered-architecture

This repository contains example code for incrementally building a layered architecture.
The code is written in Go because it's explicit yet relatively simple, but the concepts are applicable to any language.

This code is not intended to be used in production.
It is meant to demonstrate architectural concepts, not to be complete or correct.

This repository is a companion to my blog post, [Layered architecture: Creating maintainable and testable services](https://kiyanmair.com/blog/layered-architecture/).
Please read the post for explanations of the code and the concepts behind it.

# Structure

The code in this repository is organised into a single module, `layer`.
Within the module are several directories.

The `versions` directory contains the application code.
Each "version" represents the complete application at a certain stage of the abstraction process, and can be run independently of the others.
For simplicity, the code in each version is in a single package.

There is also a `sandbox` directory, whose code is equivalent to the initial version of the application (v0).
You can use the sandbox to experiment with refactoring the code yourself.
This allows you to learn hands-on, without worrying about breaking anything.
I suggest turning off AI code completion tools, like GitHub Copilot, otherwise they will likely suggest the exact same code that's elsewhere in the repository.

The `cmd` directory contains a command-line interface (CLI) for the application named `gonion`.
This enables a single `main` package to run any version of the application, including the sandbox.
It also simplifies the process of managing the database and sending requests to the endpoint.

The `db` directory contains code for interacting with the database.
The `data` directory contains the database file itself.

Finally, the `config` directory contains any shared configuration.

# Usage

To run the code, you'll need to [install Go](https://go.dev/dl/).

To build `gonion`, run `make gonion` from the root of the repository.
You will need to re-run this command if you make changes to the code.

To initialise or reset the database, run `./gonion db reset`.
To see all the available commands, run `./gonion help`.

The server accepts HTTP requests with JSON bodies, which you can send with `curl` or `gonion`.
`gonion` sends HTTP requests rather than calling the application code directly, so it's a drop-in replacement for `curl`.

You can run all tests with `make test`.

# Example

Let's register a new user.

Reset the database:

```bash
./gonion db reset
```

Start the server using the sandbox:

```bash
./gonion run --sandbox
```

Or specify a version number:

```bash
./gonion run --version N
```

Register a new user:

```bash
./gonion register --email "test@example.com"
```

View the user in the database:

```bash
./gonion db show
```
