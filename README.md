# htmx-templ-template
An end to end htmx and go-templ template using mage for deployment.

This will set you up with a working dev environment using htmx, go-templ, tailwind, and daisyUI. Note this is only configured to run a dev server. If you want to adapt this to also build production ready binaries you will have to add the mage functions yourself.

# How to use
1. Install npx in your path `npm install -g npx`
2. Install the templ command in your path `go install github.com/a-h/templ/cmd/templ@latest`
3. Install the air command in your path `go install github.com/air-verse/air@360714a021b1b77e50a5656fefc4f8bb9312d328`. NOTE if you don't install this hash you will sometimes get orphaned processes from air on mac machines. See the [air issue](https://github.com/air-verse/air/issues/534) for more details
4. Run `mage run dev`
5. Open http://localhost:7331 (the server is listening on :4433 but templ injects a watcher with an autoreload script)

When you make changes to your templ files or any of your go code everything will regenerate and then autoreload your web page