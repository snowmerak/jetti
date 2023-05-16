---
title: "Summary"
date: 2023-05-16T21:12:23+09:00
weight: 1
---

Jetti is a code generator for Golang, designed to reflect the considerations of my project structure.

## Project Structure

Jetti strives to have the following project structure:

```
.
├── cmd
├── model
├── lib
├── internal
└── doc
```

## cmd

The cmd directory contains main packages that serve as entry points for the project.

By placing multiple entry points under the cmd directory, it allows sharing of packages such as lib, model, and internal, enabling the creation of various programs.

## model

The model directory contains the models that the project will use.

It is expected to contain models of the following types:

1. protobuf or flatbuf files
2. Objects for communication (e.g., json, xml, yaml)
3. Objects for configuration (e.g., json, xml, yaml, env)

## lib

The lib directory contains the libraries that the project will use.

The package structure follows these rules:

1. The top-level package directly under lib is separated by domain.
2. Objects and actions are separated into different packages.
3. DTOs are obtained from the model package.
4. Context is used actively.

## internal

The internal directory contains packages that are used only within the project.

I'm not sure about this part.

## doc

The doc directory contains project documentation.

It houses any form of documentation that couldn't be handled as comments.

The format can be UML or a static site generator like Hugo, it doesn't matter.
