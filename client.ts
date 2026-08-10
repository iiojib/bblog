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
