const EventSource = require("eventsource");
const { stringify } = require("querystring");
const sse = new EventSource("http://localhost:3000/sse")

sse.addEventListener("notice", e => {
    const data = JSON.parse(e.data);
    if (data.Status) console.log(`notice event: ${data.Description}`)
})

sse.onmessage = e => {
    const data = JSON.parse(e.data);
    if (data.Status) console.log(`normal event: ${data.Description}`)
}

sse.onerror = e => {
    let errorMessage = String(e.message)
    if (
        errorMessage.includes("connect ECONNREFUSED") ||
        errorMessage.includes("socket hang up")
    ) {
        sse.close()
        console.log("SSE is closed")
    }
}