# htmx-templ-template
An end to end htmx and go-templ template using mage for deployment.

This will set you up with a working dev environment using htmx, go-templ, tailwind, and daisyUI. Note this is only configured to run a dev server. If you want to adapt this to also build production ready binaries you will have to add the mage functions yourself.

# How to use
1. Install mage in your path `brew install mage`. See [https://magefile.org/](https://magefile.org/) for other installation instructions
1. Run `mage install`
4. Run `mage run`
5. Open http://localhost:7331 (the server is listening on :4433 but templ injects a watcher with an autoreload script)

When you make changes to your templ files or any of your go code everything will regenerate and then autoreload your web page