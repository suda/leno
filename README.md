<div align="center">
	<img src="assets/logo.png">
	<p>
		<b>Command line log viewer with a web UI</b>
	</p>
	<br>
    <a href="https://www.npmjs.com/package/leno"><img src="https://img.shields.io/npm/v/leno" /></a>
	<a href="https://contributionswelcome.org/"><img src="https://img.shields.io/badge/contributions-welcome-7dcfef" /></a>
	<a href="https://choosealicense.com/licenses/mit/"><img src="https://img.shields.io/github/license/suda/go-gooey" /></a>
	<br>
	<br>
</div>

Leño is a [ndjson](http://ndjson.org/) log viewer with a web UI. Think of it as local Kibana/Sumo Logic for development. Works great with [Pino](https://getpino.io/) logging library or any other app that logs into JSON.

![](./assets/screenshot.png)

## 📙 Installation

```
$ npm install -g leno
```

## 🔍️ Usage

Just run your app and pipe it's output to leño:

```
$ ./myapp | leno
🌲️ Leño is running on http://localhost:3000
```

Now you can see the logs at [http://localhost:3000](http://localhost:3000).

## 🙇️ Credits

Leño's logo is based on [log by Smalllike](https://thenounproject.com/term/log/2784204) from the Noun Project.