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
