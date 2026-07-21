import client from "./client";

const defaultUrl = `http://${location?.hostname || "localhost"}:8088`;

client((prompt("SSE URL", defaultUrl) || defaultUrl).trim());
