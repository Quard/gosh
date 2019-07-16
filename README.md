# gosh

[![Build Status](https://travis-ci.com/Quard/gosh.svg?branch=master)](https://travis-ci.com/Quard/gosh)

quite simple URL shortener service without UI made in education purposes

## How to run

`go build`

`./gosh` — with a redis on localhost

`./gosh -storage bolt` — with a bolt instead of redis

## API

you can found OpenAPI v3 schema in `/api/` folder or use HTML version from `third_party/redoc-static.html`
