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
