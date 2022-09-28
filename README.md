# Ensign

**An eventing platform that is distributed in time and space.**

## Quick Start

Build the docker compose images:

```
$ ./containers/build.sh build
```

Then once the images are built run the docker containers:

```
$ ./containers/build.sh up
```

You can then run publishers and subscribers using the `debug` command. To generate 30 events per second:

```
$ go run ./cmd/debug generate
```

And to create a subscriber to consume the events:

```
$ go run ./cmd/debug consume
```

## Documentation

The primary [Ensign documentation](https://ensign.rotational.dev/) is published in this repo in the `docs/` directory as a [Hugo](https://gohugo.io/) site using the [hugo-book theme](https://github.com/alex-shpak/hugo-book).

The API documentation will be online on [pkg.go.dev](https://pkg.go.dev/) when we make this repository open source. Until then you can view the API documentation locally using [`godoc`](https://pkg.go.dev/golang.org/x/tools/cmd/godoc).

### Hugo Documentation

First, ensure you have `hugo` installed by following the [hugo installation instructions](https://gohugo.io/getting-started/installing/) for your operating system. On OS X the simplest way to do this is with homebrew:

```
$ brew install hugo
```

Change directories into the `docs/` directory. Note all of the following commands assume that your current working directory is `docs/`.

To run the local hugo development server:

```
$ hugo serve -D
```

You should now be able to view the docs locally at [http://localhost:1313](http://localhost:1313).

To create a new documentation page:

```
$ hugo new path/to/page.md
```

This will create a new page with the documentation default template; open in a browser and edit in markdown; your changes should reload live in the browser! Our theme has many options and shortcodes; to see how these work, please viist the theme documentation. A couple of important quick links are below:

- [Page Configuration](https://github.com/alex-shpak/hugo-book#page-configuration)
- [Hints and Notification Blocks](https://hugo-book-demo.netlify.app/docs/shortcodes/hints/)
- [Mermaid Diagrams](https://hugo-book-demo.netlify.app/docs/shortcodes/mermaid/)
- [KaTeX math typesetting](https://hugo-book-demo.netlify.app/docs/shortcodes/katex/)
- [Columns](https://hugo-book-demo.netlify.app/docs/shortcodes/columns/) and [Tabs](https://hugo-book-demo.netlify.app/docs/shortcodes/tabs/)

To help with generating and updating markdown tables, we recommend:

- [Markdown Table Generator](https://www.tablesgenerator.com/markdown_tables)

### API Documentation

First, ensure you have `godoc` installed:

```
$ go install golang.org/x/tools/cmd/godoc@latest
```

You can then run the API documentation locally using:

```
$ godoc --http=:6060
```

The docs should be available shortly when you open your browser to [http://localhost:6060/pkg/github.com/rotationalio/ensign/pkg/](http://localhost:6060/pkg/github.com/rotationalio/ensign/pkg/).