---
title: "Developer SDKs"
weight: 20
date: 2023-05-17T17:03:41-04:00
---

Welcome Ensign user!

If you're looking to figure out how to write code to connect to Ensign, you've come to the right place.

<!--more-->

Ensign is written primarily in Go and Protocol Buffers because that's what we like, but hey, this isn't about us, it's about **you**.

We want Ensign to be as accessible as possible, so we made you some SDKs 💙

## Available SDKs

Here is a list of the SDKs that are currently available.

You should be able to install each using your usual language-specific package manager (e.g. `go get`, `pip install`, `npm install`, etc.):

- [go-ensign](https://github.com/rotationalio/go-ensign)
- [pyensign](https://github.com/rotationalio/pyensign)
- [ensignJS](https://github.com/rotationalio/ensignjs) *Note: 🙈 this is just an empty repo for now, but we're working on it!!*

Each SDK is structured in a language-specific fashion and contains cross-sdk driver methods as well as tooling that may be tailored to users of specific languages.

Don't see an Ensign SDK for the language you love most? Tell us (support@rotational.io), and we'll get right on it!

## What's in the SDK Repos?

The SDK repos above contain driver and library code for interacting with Ensign. They each have their own language-specific docs.

### Common Methods

Here are some of the common driver methods that all of the SDKs have in some shape or form (*there are some implementation differences due to the ways that different languages handle concurrency*).

- Create an open `Client` to Ensign using your [API keys]({{< ref "/getting-started#getting-started" >}})
- Create a new `Topic` on your `Client`
- Structure valid Ensign `Events` of various possible `MimeTypes`
- Create a new `Publisher` on your `Client` and invoke this `Publisher`'s `Publish` method to publish an `Event` to a `Topic`
- Create a new `Subscriber` on your `Client` and create an event stream to collect published `Events` using this `Subscriber`'s `Subscribe` method

## What are these Protobuf Things?

Protocol buffers are useful and help us make Ensign real fast, but we know they aren't exactly common.

If you're familiar with other serialization formats like JSON and XML and are just getting up to speed with protobuf, check out [this post](https://rotational.io/blog/what-are-protocol-buffers/).

If you're less familiar with serialization methods, the main thing to understand is that the protocol buffers are what define the Ensign *service*; they're a set of eventing-related rules and components that Ensign understands, and if you want to write code that interacts with Ensign, you have to explain what you want in terms of those rules and components.

Protocol buffers are like a recipe for ingredients that have to be combined and baked (aka compiled) into code like a casserole before you can ~~eat~~ use it. But compiling protocol buffers is not always the easiest thing, so our SDK libraries each have a copy of the protobufs compiled in the SDK language.

You shouldn't (hopefully) have to jump through the compilation hoops &mdash; just install and import the SDK you need in the language you prefer. Hopefully this makes it as easy as possible for you!
