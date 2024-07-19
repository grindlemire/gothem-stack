# gothem-stack
An end to end htmx and go-templ template using echo for the web server and mage for deployment.

**Go**\
**T**empl\
**H**TMX\
**E**cho\
**M**age

For the frontend libraries it uses tailwindcss and daisyUI for styling. You can easily integrate alpine-js 

This will set you up with a working dev environment using htmx, go-templ, tailwind, and daisyUI.

# How to use
1. Install mage in your path `brew install mage`. See [https://magefile.org/](https://magefile.org/) for other installation instructions
1. Run `mage install`
4. Run `mage run`
5. Open http://localhost:7331 (the server is listening on :4433 but templ injects a watcher with an autoreload script)

When you make changes to your templ files or any of your go code everything will regenerate and then autoreload your web page

# Other commands
`mage templ` - Do a one time regeneration of your templ files

`mage build` - Do a one time build of the go files

`mage tidy` - Run `go mod tidy`