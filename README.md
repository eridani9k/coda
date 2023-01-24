# Routing API (Round-Robin)

## Introduction

This repository contains code for an example routing API backed by a round-robin load balancing algorithm. The code was designed to be bare-bones in terms of setup and infrastructure for the sole purpose of code quality review.

A full deployment of this example consists of:
- 1 `Router` which acts as a _reverse proxy_
- 1..N backend `API` servers load balanced by the `Router`.

## Usage


