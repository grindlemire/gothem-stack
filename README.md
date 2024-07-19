<div align="center" style="text-align: center;">
  <img src="gopher-batman.png" alt="Description" style="width: 50%; max-width: 300px;">
</div>

<p align="center"><a href="https://pkg.go.dev/github.com/Permify/policy-enforcer?tab=doc" 
target="_blank"></a><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go" alt="go version" />&nbsp;&nbsp;<img src="https://img.shields.io/github/license/grindlemire/gothem-stack?style=for-the-badge" alt="license" />

# gothem-stack

An end to end [htmx](https://htmx.org) and [go-templ](https://templ.guide) template using [echo](https://echo.labstack.com/) for the web server and [mage](https://magefile.org/) for deployment. Other than that this is unopinionated so bring all your other favorite technologies :)

**Go**\
**T**empl\
**H**TMX\
**E**cho\
**M**age

For the frontend libraries it uses [tailwindcss](https://tailwindcss.com/) and [daisyUI](https://daisyui.com/) for styling. You can easily integrate [alpine-js](https://alpinejs.dev/) if you desire to.

This will set you up with a working dev environment using htmx, go-templ, tailwind, and daisyUI.

## Quickstart
1. Install npm in your path `brew install node; brew install npm`
1. Install mage in your path `brew install mage`. See [https://magefile.org/](https://magefile.org/) for other installation instructions
1. Run `mage install`
1. Run `mage run`
1. Open http://localhost:7331 (the server is listening on :4433 but templ injects a watcher with an autoreload script)

When you make changes to your templ files or any of your go code everything will regenerate and then autoreload your web page

## How to Use
`mage run` - Run an interactive development environment that will automatically reload on any file change. Listens on port :4433 and has an autoreload page on :7331

`mage install` - Install all the dependencies

`mage templ` - Do a one time regeneration of your templ files

`mage build` - Do a one time build of the go files

`mage tidy` - Run `go mod tidy`

## What about persistence?
I usually use postgres and [sqlc](https://sqlc.dev/) and they work great. If people start to use this and want to see that added in I'll do so.

### But Gotham is spelled with an 'a'....
Yea I know it's spelled with an 'a' :] I was trying to come up with a name that was easier to say and interact with than 'htmx-templ-template' and gothem is what I came up with.