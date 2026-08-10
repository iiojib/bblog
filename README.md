# bblog

Tiny stdin-to-browser log streaming tool.

## Features

- Works without code changes.
- Activates via bookmarklet.
- Translates ANSI SGR escape codes into browser console styles.

## Install

### Using installer script

```bash
curl -fsSL https://github.com/iiojib/bblog/releases/latest/download/install.sh | sh
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
javascript:(()%3D%3E%7Bfunction%20e(r)%7Blet%20o%3Dconsole%2Ct%3Dnew%20EventSource(r)%3Bt.onopen%3D()%3D%3Eo.log(%60%5BbbLog%5D%20connected%3A%20%24%7Br%7D%60)%2Ct.onmessage%3Dn%3D%3E%7Blet%20c%3DString(n.data)%3Bif(c%3D%3D%3D%60%0A__BBLOG_SHUTDOWN__%60)return%20o.log(%22%5BbbLog%5D%20server%20shutdown%22)%2Ct.close()%3Bo.log(...c.split(%60%0A%60))%7D%2Ct.onerror%3Dn%3D%3E%7Bif(t.readyState%3D%3D%3D0)return%20o.log(%22%5BbbLog%5D%20connection%20lost%2C%20reconnecting...%22)%3Bo.error(%22error%3A%22%2Cn)%7D%7Dvar%20s%3D%60http%3A%2F%2F%24%7Blocation%3F.hostname%7C%7C%22localhost%22%7D%3A8088%60%3Be((prompt(%22SSE%20URL%22%2Cs)%7C%7Cs).trim())%3B%7D)()%3B
```

Open your app in the browser, then click the bookmarklet to start streaming logs.

## Options

- `[path]` read from named pipe (created if absent), or stdin if omitted
- `-H string` HTTP listen host (default: `0.0.0.0`)
- `-P int` HTTP listen port (default: `8088`)
- `-S` strip all ANSI escape codes (emit plain text)
- `-max-buffer-size int` Maximum buffer size in KB, minimum: 4KB (default: `64`)
- `-v` print version and exit

## Advanced Usage

Many apps disable color output when they detect piped stdout. You can force color output by setting the `FORCE_COLOR` environment variable:

```bash
FORCE_COLOR=1 myapp 2>&1 | tee -i >(bblog)
```

---

**Some apps may also buffer output when they detect piped stdout. Here are a few ways to disable buffering:**

- Python and `sed` support the `-u` flag to disable buffering (see the `sed` example below).
- Other programs may support flags like `--line-buffered`, `--unbuffered`, etc.
- For Python, you can also set `PYTHONUNBUFFERED=1`.
- Apps that rely on C stdio may support the `stdbuf` command to disable buffering.
- Or you can use the `unbuffer` command from the `expect` package to run your application in a PTY.

---

You can add a prefix to the log output to help identify which app the logs are coming from:

```bash
myapp 2>&1 | sed -u 's/^/[myapp] /' | tee -i >(bblog)
```

Or apply the prefix only to the broadcasted stream:

```bash
myapp 2>&1 | tee -i >(sed -u 's/^/[myapp] /' | bblog)
```

`tee -i` already ignores SIGINT. If you want to preserve shutdown logs from your app while using additional processing commands like `sed`, wrap those commands in `(trap '' INT; ...)`:

```bash
myapp 2>&1 | (trap '' INT; sed -u 's/^/[myapp] /') | tee -i >(bblog)

# or

myapp 2>&1 | tee -i >((trap '' INT; sed -u 's/^/[myapp] /') | bblog)
```

---

You can also add a timestamp to each line using `awk`:

```bash
myapp 2>&1 | awk '{ print strftime("[%H:%M:%S]"), $0; fflush(); }' | tee -i >(bblog)
```

---

If you want the browser connection to survive manual app restarts, start `bblog` with a named pipe:

```bash
bblog /tmp/myapp.pipe
```

Then run your app and pipe its output to that named pipe:

```bash
myapp 2>&1 | tee -i /tmp/myapp.pipe
```

With this approach you can also stream logs from multiple apps by piping their outputs to the same named pipe.

---

Add the following snippet to your HTML page to connect automatically to the log stream:

```html
<script>
	(()=>{function t(r){let o=console,e=new EventSource(r);e.onopen=()=>o.log(`[bbLog] connected: ${r}`),e.onmessage=n=>{let c=String(n.data);if(c===`
__BBLOG_SHUTDOWN__`)return o.log("[bbLog] server shutdown"),e.close();o.log(...c.split(`
`))},e.onerror=n=>{if(e.readyState===0)return o.log("[bbLog] connection lost, reconnecting...");o.error("error:",n)}}globalThis.bbLogClient=t;})();

  bbLogClient("http://localhost:8088");
</script>
```

Or use this TypeScript snippet:

```typescript
export default function client(url: string): void {
	const c = console;
	const sse = new EventSource(url);

	sse.onopen = () => c.log(`[bbLog] connected: ${url}`);
	sse.onmessage = (event) => {
		const data = String(event.data);

		if (data === "\n__BBLOG_SHUTDOWN__") {
			c.log("[bbLog] server shutdown");
			return sse.close();
		}

		c.log(...data.split("\n"));
	};

	sse.onerror = (event) => {
		if (sse.readyState === 0) return c.log("[bbLog] connection lost, reconnecting...");
		c.error("error:", event);
	};
}
```
