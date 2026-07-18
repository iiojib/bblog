import client from "./client";

const defaultHost = location?.hostname || "localhost";
const defaultHostPort = `${defaultHost}:8088`;

client((prompt("SSE host:port", defaultHostPort) || defaultHostPort).trim());
