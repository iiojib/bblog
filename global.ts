import client from "./client";

declare global {
	var bbLogClient: typeof client;
}

globalThis.bbLogClient = client;
