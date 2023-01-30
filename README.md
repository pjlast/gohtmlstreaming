# HTML Streaming

The purpose of this repo is to demonstrate the potential benefits of HTML
streaming.

It's implemented in Go simply because Go has a straightforward, batteries included
approach when it comes to HTTP servers, and a powerful HTML templating library
that makes it extremely easy to modularise streamable parts of the page.

There are two pages: `/stream` and `/nostream`. `/nostream` functions in the
way that most web pages are implemented by default: the page is only rendered
once all content has been received.

A small loading element is streamed just before waiting for fetched content to
arrive, and an inline CSS style is streamed to hide that loading element again
after the content is ready to be streamed.

## Run

You'll need a working version of Go installed:

```
go run .
```
