# bblog

Tiny stdin-to-browser log streaming tool.

## Features

- Works without code changes.
- Activates through a bookmarklet.
- Translates ANSI escape codes into browser console styles.

## Install

### Using installer script

```bash
curl -L https://github.com/iiojib/bblog/releases/latest/download/install.sh | sh
```

### From source

```bash
git clone https://github.com/iiojib/bblog.git
cd bblog
go install .
```

## Usage

Start your app and pipe its output to `bblog`:

```bash
myapp 2>&1 | tee -i >(bblog)
```

Save the following bookmarklet as a browser bookmark URL:

```bash
javascript:(()%3D%3E%7Bfunction%20e(t)%7Blet%20o%3Dnew%20EventSource(t)%3Bo.onopen%3D()%3D%3Econsole.log(%60%5BbbLog%5D%20connected%3A%20%24%7Bt%7D%60)%2Co.onmessage%3Dn%3D%3E%7Blet%20r%3DString(n.data)%3Bif(r%3D%3D%3D%60%0A__BBLOG_SHUTDOWN__%60)return%20console.log(%22%5BbbLog%5D%20server%20shutdown%22)%2Co.close()%3Bconsole.log(...r.split(%60%0A%60))%7D%2Co.onerror%3Dn%3D%3E%7Bif(o.readyState%3D%3D%3DEventSource.CONNECTING)return%20console.log(%22%5BbbLog%5D%20connection%20lost%2C%20reconnecting...%22)%3Bconsole.error(%22error%3A%22%2Cn)%7D%7Dvar%20c%3D%60http%3A%2F%2F%24%7Blocation%3F.hostname%7C%7C%22localhost%22%7D%3A8088%60%3Be((prompt(%22SSE%20URL%22%2Cc)%7C%7Cc).trim())%3B%7D)()%3B
```

Open your app in the browser, then click the bookmarklet to start streaming logs.

## CLI flags:

- `-H string` HTTP listen host (default: `0.0.0.0`)
- `-P int` HTTP listen port (default: `8088`)
- `-S` strip ANSI escape codes (emit plain text)

## Advanced Usage

Many apps disable color output when they detect piped stdout. You can force color output by setting the `FORCE_COLOR` environment variable:

```bash
FORCE_COLOR=1 myapp 2>&1 | tee -i >(bblog)
```

---

If you want the connection to survive manual restarts of your app, pipe the output to a log file:

```bash
myapp 2>&1 | tee -i -a /tmp/myapp.log
```

Then run `bblog` in another terminal session to stream that log file:

```bash
tail -f /tmp/myapp.log | bblog
```

With this approach you can also stream logs from multiple apps by piping their outputs to the same log file.

---

Add the following snippet to your HTML page to connect automatically to the log stream:

```html
<script>
  (()=>{function n(t){let o=new EventSource(t);o.onopen=()=>console.log(`[bbLog] connected: ${t}`),o.onmessage=e=>{let r=String(e.data);if(r===`
__BBLOG_SHUTDOWN__`)return console.log("[bbLog] server shutdown"),o.close();console.log(...r.split(`
`))},o.onerror=e=>{if(o.readyState===EventSource.CONNECTING)return console.log("[bbLog] connection lost, reconnecting...");console.error("error:",e)}}globalThis.bbLogClient=n;})();

  bbLogClient("http://localhost:8088");
</script>
```

Or use this TypeScript snippet:

```typescript
export default function client(url: string): void {
	const sse = new EventSource(url);

	sse.onopen = () => console.log(`[bbLog] connected: ${url}`);
	sse.onmessage = (event) => {
		const data = String(event.data);

		if (data === "\n__BBLOG_SHUTDOWN__") {
			console.log("[bbLog] server shutdown");
			return sse.close();
		}

		console.log(...data.split("\n"));
	};

	sse.onerror = (event) => {
		if (sse.readyState === EventSource.CONNECTING) return console.log("[bbLog] connection lost, reconnecting...");
		console.error("error:", event);
	};
}
```
