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
myapp 2>&1 | tee >(bblog)
```

Save the following bookmarklet as a browser bookmark URL:

```bash
javascript:(()%3D%3E%7Bfunction%20n(c)%7Blet%20l%3D%60%0A__BBLOG_SHUTDOWN__%60%2Ce%3D%60http%3A%2F%2F%24%7Bc%7D%2F%60%2Co%3Dnew%20EventSource(e)%3Bo.onopen%3D()%3D%3Econsole.log(%60%5BbbLog%5D%20connected%3A%20%24%7Be%7D%60)%2Co.onmessage%3Dt%3D%3E%7Blet%20r%3DString(t.data)%3Bif(r%3D%3D%3Dl)return%20console.log(%22%5BbbLog%5D%20server%20shutdown%22)%2Co.close()%3Bconsole.log(...r.split(%60%0A%60))%7D%2Co.onerror%3Dt%3D%3E%7Bif(o.readyState%3D%3D%3DEventSource.CONNECTING)return%20console.log(%22%5BbbLog%5D%20connection%20lost%2C%20reconnecting...%22)%3Bconsole.error(%22error%3A%22%2Ct)%7D%7Dvar%20i%3Dlocation%3F.hostname%7C%7C%22localhost%22%2Cs%3D%60%24%7Bi%7D%3A8088%60%3Bn((prompt(%22SSE%20host%3Aport%22%2Cs)%7C%7Cs).trim())%3B%7D)()%3B
```

Open your app in the browser, then click the bookmarklet to start streaming logs.

## Options and Arguments

CLI flags:

- `-H string` HTTP listen host (default: `0.0.0.0`)
- `-P int` HTTP listen port (default: `8088`)
- `-N` disable timestamp in emitted messages
- `-S` strip ANSI escape codes (emit plain text)

Positional arguments:

- `prefix` optional. Appended as `[prefix] ` before each log line.

## Advanced Usage

Many apps disable color output when they detect piped stdout. You can force color output by setting the `FORCE_COLOR` environment variable:

```bash
FORCE_COLOR=1 myapp 2>&1 | tee >(bblog)
```

Add the following snippet to your HTML page to connect automatically to the log stream:

```html
<script>
  (()=>{function n(l){let c=`
__BBLOG_SHUTDOWN__`,t=`http://${l}/`,o=new EventSource(t);o.onopen=()=>console.log(`[bbLog] connected: ${t}`),o.onmessage=e=>{let r=String(e.data);if(r===c)return console.log("[bbLog] server shutdown"),o.close();console.log(...r.split(`
`))},o.onerror=e=>{if(o.readyState===EventSource.CONNECTING)return console.log("[bbLog] connection lost, reconnecting...");console.error("error:",e)}}globalThis.bbLogClient=n;})();

  bbLogClient("localhost:8088");
</script>
```

Or use this TypeScript snippet:

```typescript
export default function client(hostPort: string): void {
  const shutdownSentinel = "\n__BBLOG_SHUTDOWN__";
  const url = `http://${hostPort}/`;
  const sse = new EventSource(url);

  sse.onopen = () => console.log(`[bbLog] connected: ${url}`);
  sse.onmessage = (event) => {
    const data = String(event.data);

    if (data === shutdownSentinel) {
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
